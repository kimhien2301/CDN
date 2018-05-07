#!/usr/bin/env ruby
# coding: utf-8

require 'rexml/document'
require 'json'
require 'optparse'

# デフォルト設定一覧
cache_params = {
  random:      { Capacity: 100                      },
  fifo:        { Capacity: 100                      },
  srrip:       { Capacity: 100, RRPVbit: 2          },
  arc:         { Capacity: 100                      },
  lirs:        { Capacity: 100                      },
  admission:   { AdmissionList: %w(1..100).to_a     },
  lfu:         { Capacity: 100                      },
  iclfu:       { Capacity: 100                      },
  windowlfu:   { Capacity: 100, Window: 1000        },
  modifiedlru: { Capacity: 100, Jump: 5             },
  lru:         { Capacity: 100                      },
  lruk:        { Capacity: 100, K: 5                },
  iris:        { Capacity: 100, SpectrumRatio: 0.95 },
}

defaults = {
  NetworkId: nil,
  RequestModels: { Gamma: { K: 0.475, Theta: 170.6067 }, Zipf: { Skewness: 0.679 }},
  OriginServer: { LibrarySize: 1000 }, CacheServers: { CacheAlgorithm: 'iris' },
  Links: { Cost: 1.0, Bidirectional: true },
  Clients: { RequestModelId: 'gamma', TrafficWeight: 1.0 }
}

# デフォルトの設定は引数で上書きできる
opt = OptionParser.new
opt.on('-f', '--graphml-file=VAL')    { |value| defaults[:NetworkId]                       = value       }
opt.on('-k', '--gamma-k=VAL')         { |value| defaults[:RequestModels][:Gamma][:K]       = value.to_f  }
opt.on('-t', '--gamma-theta=VAL')     { |value| defaults[:RequestModels][:Gamma][:Theta]   = value.to_f  }
opt.on('-s', '--zipf-s=VAL')          { |value| defaults[:RequestModels][:Zipf][:Skewness] = value.to_f  }
opt.on('-l', '--library-size=VAL')    { |value| defaults[:OriginServer][:LibrarySize]      = value.to_i  }
opt.on('--[no-]bidirectional')        { |value| defaults[:Links][:Bidirectional]           = value       }
opt.on('-r', '--request-model=VAL')   { |value| defaults[:Clients][:RequestModelId]        = value.to_s  }
opt.on('-a', '--cache-algorithm=VAL') { |value|
  if cache_params.keys.map(&:to_s).include?(value)
    defaults[:CacheServers][:CacheAlgorithm] = value
  else
    puts "error: invalid cache algorithm"
    puts "available algorithm: " + cache_params.keys.map(&:to_s).join(", ")
    exit 1
  end
}
opt.on('-C', '--cache-capacity=VAL') { |value|
  cache_params.keys.each do |algorithm|
    if cache_params[algorithm][:Capacity]
      cache_params[algorithm][:Capacity] = value.to_i
    end
  end
}
opt.on('-B', '--RRPVbit=VAL')       { |value| cache_params[:srrip][:RRPVbit]       = value.to_i }
opt.on('-K', '--Kth=VAL')           { |value| cache_params[:lruk][:K]              = value.to_i }
opt.on('-J', '--Jump=VAL')          { |value| cache_params[:modifiedlru][:Jump]    = value.to_i }
opt.on('-W', '--Window=VAL')        { |value| cache_params[:windowlfu][:Window]    = value.to_i }
opt.on('-R', '--SpectrumRatio=VAL') { |value| cache_params[:iris][:SpectrumRatio]  = value.to_f }
opt.parse!(ARGV)

# 読み込むファイル名が指定されていない場合は使い方を表示して終了
unless defaults[:NetworkId]
  puts "Please specify input graphml file with -f option."
  puts opt.help
  exit 1
end

# ファイルを読み込んでXMLをパース
doc = REXML::Document.new(open(defaults[:NetworkId]))

# ネットワーク名を拡張子を除いたファイル名に設定
defaults[:NetworkId] = File.basename(defaults[:NetworkId]).split('.')[0..-2].join('.')

# JSONファイルの元となるHashを作成
graph = {}
graph[:NetworkId] = defaults[:NetworkId]
graph[:RequestModels] = [
  { Id: "gamma", Model: "gamma", ParameterKeys: ["K", "Theta"], Parameters: [defaults[:RequestModels][:Gamma][:K], defaults[:RequestModels][:Gamma][:Theta]] },
  { Id: "zipf", Model: "zipf", ParameterKeys: ["Skewness"], Parameters: [defaults[:RequestModels][:Zipf][:Skewness]] }
]

graph[:OriginServer] = nil
graph[:CacheServers] = []
graph[:Links] = []
graph[:Clients] = []

# graphmlファイルからノード情報を読み込む
origin_id = nil
doc.elements.each('graphml/graph/node') do |node_element|
  id = node_element.attributes['id']
  label = node_element.elements.each('data/y:ShapeNode/y:NodeLabel'){}.first.text

  if label == 'origin'
    # ラベルが "origin" に設定されているノードをオリジンサーバとして扱う
    origin_id = id
    graph[:OriginServer] = { Id: id, LibrarySize: defaults[:OriginServer][:LibrarySize] }
  else
    # それ以外のノードはキャッシュサーバとして扱う
    graph[:CacheServers] << { Id: id, CacheAlgorithm: defaults[:CacheServers][:CacheAlgorithm], ParameterKeys:  cache_params[defaults[:CacheServers][:CacheAlgorithm].to_sym].keys.map(&:to_s), Parameters:     cache_params[defaults[:CacheServers][:CacheAlgorithm].to_sym].values }
  end
end

# graphmlファイルからリンク情報を読み込む
origin_connected = false
doc.elements.each('graphml/graph/edge') do |edge_element|
  id = edge_element.attributes['id']
  source = edge_element.attributes['source']
  target = edge_element.attributes['target']

  if (source == origin_id and target != origin_id) or (source != origin_id and target == origin_id)
    origin_connected = true
  end

  graph[:Links] << { EdgeNodeIds:   [source, target], Cost:          defaults[:Links][:Cost], Bidirectional: defaults[:Links][:Bidirectional] }
end

if origin_id.nil?
  # オリジンサーバが存在しないものは不正なグラフとして扱う
  puts 'Error: This graphml file does not contain a node labeled "origin".'
  exit 1
elsif origin_connected == false
  # オリジンサーバの接続リンクがないものはファイルを取得できないため不正なグラフとして扱う
  puts 'Error: Origin node does not connected to any other node.'
  exit 1
end

# 各キャッシュサーバの下にクライアントを接続する
graph[:CacheServers].each do |node|
  graph[:Clients] << { UpstreamId: node[:Id], RequestModelId: defaults[:Clients][:RequestModelId], TrafficWeight:  defaults[:Clients][:TrafficWeight] }
end

# JSON形式に整形して出力
puts JSON.pretty_generate(graph)
