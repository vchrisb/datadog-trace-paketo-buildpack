/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package trace_test

import (
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/datadog-trace/trace"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect trace.Detect
	)

	it("fails without service and variable", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
	})

	it("passes with service", func() {
		ctx.Platform.Bindings = libcnb.Bindings{
			{Name: "test-service", Type: "DatadogTrace"},
		}

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "datadog-trace-java"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "datadog-trace-java"},
						{Name: "jvm-application"},
					},
				},
			},
		}))
	})

	context("$BP_DATADOGTRACE_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_DATADOGTRACE_ENABLED", "true")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_DATADOGTRACE_ENABLED")).To(Succeed())
		})

		it("passes with BP_DATADOGTRACE_ENABLED", func() {
			Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
				Pass: true,
				Plans: []libcnb.BuildPlan{
					{
						Provides: []libcnb.BuildPlanProvide{
							{Name: "datadog-trace-java"},
						},
						Requires: []libcnb.BuildPlanRequire{
							{Name: "datadog-trace-java"},
							{Name: "jvm-application"},
						},
					},
				},
			}))
		})
	})
}
