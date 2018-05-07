package main

import "fmt"
import "utils"
import "encoding/json"

type NodeReport map[string]interface{}
type LinkReport map[string]interface{}

func (g *Graph_t) PlainReport() {
	utils.Print(fmt.Sprintln("Nodes:"))
	for index, node := range g.nodes {
		hit := node.Entity().(*ServerModel_t).Storage().HitCount()
		miss := node.Entity().(*ServerModel_t).Storage().MissCount()
		utils.Print(fmt.Sprintf("  [%3d]\t%s\t(access:%6d,\thit:%6d,\thit rate:%6.1f%%)\n", index, node.ID(), hit+miss, hit, float64(hit)/float64(hit+miss)*100.0))
	}
	utils.Print(fmt.Sprintln("Links:"))
	for index, link := range g.links {
		utils.Print(fmt.Sprintf("  [%3d]\t%s -> %s,\ttraffic:%6.1f\n", index, link.src.ID(), link.dst.ID(), link.traffic))
	}
}

func (g *Graph_t) JsonReport() {
	report := make(map[string]interface{})
	summary := make(map[string]interface{})
	originServerReport := make(map[string]interface{})
	cacheServerReports := make([]NodeReport, 0)
	linkReports := make([]LinkReport, 0)

	summary["internalTraffic"] = 0.0
	summary["originTraffic"] = 0.0

	origin := g.originServers()[0]
	originServerReport["id"] = origin.ID
	originServerReport["hit"] = origin.Entity().(*ServerModel_t).Storage().HitCount()
	originServerReport["miss"] = origin.Entity().(*ServerModel_t).Storage().MissCount()
	originServerReport["accesses"] = originServerReport["hit"].(int) + originServerReport["miss"].(int)
	originServerReport["hitrate"] = origin.Entity().(*ServerModel_t).hitRate()
	originServerReport["caches"] = origin.Entity().(*ServerModel_t).Storage().CacheList()

	for _, node := range g.cacheServers() {
		nodeReport := make(map[string]interface{})
		nodeReport["id"] = node.ID
		nodeReport["hit"] = node.Entity().(*ServerModel_t).Storage().HitCount()
		nodeReport["miss"] = node.Entity().(*ServerModel_t).Storage().MissCount()
		nodeReport["accesses"] = nodeReport["hit"].(int) + nodeReport["miss"].(int)
		nodeReport["hitrate"] = node.Entity().(*ServerModel_t).hitRate()
		nodeReport["caches"] = node.Entity().(*ServerModel_t).Storage().CacheList()
		cacheServerReports = append(cacheServerReports, nodeReport)
	}
	for _, link := range g.links {
		linkReport := make(map[string]interface{})
		linkReport["src"] = link.src.ID()
		linkReport["dst"] = link.dst.ID()
		linkReport["traffic"] = link.traffic

		if link.src.ID() == originServerReport["id"] || link.dst.ID() == originServerReport["id"] {
			summary["originTraffic"] = summary["originTraffic"].(float64) + link.traffic
		} else {
			summary["internalTraffic"] = summary["internalTraffic"].(float64) + link.traffic
		}

		linkReports = append(linkReports, linkReport)
	}
	summary["totalTraffic"] = summary["internalTraffic"].(float64) + summary["originTraffic"].(float64)

	report["OriginServer"] = originServerReport
	report["CacheServers"] = cacheServerReports
	report["Links"] = linkReports
	jsonString, _ := json.Marshal(report)
	utils.Print(fmt.Sprintln(string(jsonString)))
}
