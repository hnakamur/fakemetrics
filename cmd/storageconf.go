// Copyright © 2018 Grafana Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"time"

	"github.com/raintank/fakemetrics/out"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/raintank/schema.v1"
)

// storageconfCmd represents the storageconf command
var storageconfCmd = &cobra.Command{
	Use:   "storageconf",
	Short: "Sends out a small set of metrics which you can test aggregation and retention rules on",
	Run: func(cmd *cobra.Command, args []string) {
		initStats(true, "storageconf")
		period = int(periodDur.Seconds())
		flush = int(flushDur.Nanoseconds() / 1000 / 1000)
		outs := getOutputs()
		if len(outs) == 0 {
			log.Fatal("need to define an output")
		}
		to := time.Now().Unix()
		from := to - 24*60*60

		names := []string{
			"fakemetrics.raw.min",
			"fakemetrics.raw.max",
			"fakemetrics.raw.sum",
			"fakemetrics.raw.default",
			"fakemetrics.agg.min",
			"fakemetrics.agg.max",
			"fakemetrics.agg.sum",
			"fakemetrics.agg.default",
		}
		var metrics []*schema.MetricData
		for _, name := range names {
			metrics = append(metrics, buildMetric(name, 1))
		}
		do(metrics, from, to, outs)
	},
}

func init() {
	rootCmd.AddCommand(storageconfCmd)
}

func buildMetric(name string, org int) *schema.MetricData {
	out := &schema.MetricData{
		Name:     name,
		Metric:   name,
		OrgId:    org,
		Interval: 1,
		Unit:     "ms",
		Mtype:    "gauge",
		Tags:     nil,
	}
	out.SetId()
	return out
}
func do(metrics []*schema.MetricData, from, to int64, outs []out.Out) {
	for ts := from; ts < to; ts++ {
		if ts%3600 == 0 {
			fmt.Println("doing ts", ts)
		}
		for i := range metrics {
			metrics[i].Time = int64(ts)
			metrics[i].Value = float64(ts % 10)
		}
		for _, out := range outs {
			err := out.Flush(metrics)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}
