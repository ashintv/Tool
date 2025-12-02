import React, { useState, useEffect } from "react";
import { useParams, useNavigate, useSearchParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowLeft, Server, Box, BarChart3, Play, Square, RotateCcw } from "lucide-react";
import type { Machine, Container, Metrics } from "@/types/types";
import axios from "axios";

type ActiveView = "machine" | "containers" | "metrics";

interface ContainerAction {
  id: string;
  action: "start" | "stop" | "restart";
}

export const MachinePage = () => {
  const { machineId } = useParams<{ machineId: string }>();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const initialView = (searchParams.get("view") as ActiveView) || "machine";
  const [activeView, setActiveView] = useState<ActiveView>(initialView);
  const [machine, setMachine] = useState<Machine | null>(null);
  const [containers, setContainers] = useState<Container[]>([]);
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [containerActions, setContainerActions] = useState<ContainerAction[]>([]);

  useEffect(() => {
    if (machineId) {
      fetchMachineData();
    }
  }, [machineId]);

  useEffect(() => {
    // Set initial view from URL parameters
    const viewFromUrl = searchParams.get("view") as ActiveView;
    if (viewFromUrl && ["machine", "containers", "metrics"].includes(viewFromUrl)) {
      setActiveView(viewFromUrl);
    }
  }, [searchParams]);

  useEffect(() => {
    // Update URL when view changes
    const newSearchParams = new URLSearchParams();
    if (activeView !== "machine") {
      newSearchParams.set("view", activeView);
    }
    const newUrl = `/machine/${machineId}${newSearchParams.toString() ? `?${newSearchParams.toString()}` : ""}`;
    window.history.replaceState({}, "", newUrl);
  }, [activeView, machineId]);

  const fetchMachineData = async () => {
    try {
      setIsLoading(true);
      // Fetch machine details
      // const machineRes = await axios.get(`http://localhost:8080/api/machine/${machineId}`, {
      //   headers: { Authorization: localStorage.getItem("token") },
      // });

      // Mock machine data for now
      setMachine({
        id: machineId!,
        name: `Machine ${machineId}`,
        status: "online",
        IP: "192.168.1.100",
        os: "Ubuntu 22.04",
        cpu: "Intel i7-12700K",
        memory: "32GB",
        storage: "1TB SSD",
        uptime: "15 days, 4 hours",
        users: ["admin", "developer", "monitor"]
      });

      // Mock containers data
      setContainers([
        { id: "nginx-1", name: "nginx-web", status: "running", image: "nginx:latest", ports: "80:80" },
        { id: "postgres-1", name: "database", status: "running", image: "postgres:15", ports: "5432:5432" },
        { id: "redis-1", name: "cache", status: "stopped", image: "redis:7", ports: "6379:6379" },
        { id: "app-1", name: "web-app", status: "running", image: "node:18", ports: "3000:3000" },
      ]);

      // Mock metrics data
      setMetrics({
        cpu: { current: 45, average: 38, peak: 82 },
        memory: { current: 68, average: 55, peak: 89 },
        network: { upload: 125, download: 450 },
        uptime: 99.8
      });
    } catch (error) {
      console.error("Error fetching machine data:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleContainerAction = async (containerId: string, action: "start" | "stop" | "restart") => {
    setContainerActions(prev => [...prev, { id: containerId, action }]);

    try {
      // Mock API call
      await new Promise(resolve => setTimeout(resolve, 1000));

      // Update container status
      setContainers(prev => prev.map(container =>
        container.id === containerId
          ? { ...container, status: action === "stop" ? "stopped" : "running" }
          : container
      ));
    } catch (error) {
      console.error(`Error ${action}ing container:`, error);
    } finally {
      setContainerActions(prev => prev.filter(item => item.id !== containerId));
    }
  };

  const renderMachineView = () => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Basic Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div><span className="font-medium">Name:</span> {machine?.name}</div>
            <div><span className="font-medium">IP Address:</span> {machine?.IP}</div>
            <div><span className="font-medium">Status:</span>
              <span className={`ml-2 px-2 py-1 rounded text-sm ${
                machine?.status === "online"
                  ? "bg-green-100 text-green-700"
                  : "bg-red-100 text-red-700"
              }`}>
                {machine?.status}
              </span>
            </div>
            <div><span className="font-medium">Uptime:</span> {machine?.uptime}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">System Specifications</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div><span className="font-medium">OS:</span> {machine?.os}</div>
            <div><span className="font-medium">CPU:</span> {machine?.cpu}</div>
            <div><span className="font-medium">Memory:</span> {machine?.memory}</div>
            <div><span className="font-medium">Storage:</span> {machine?.storage}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Users & Access</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <span className="font-medium">Authorized Users:</span>
              <div className="flex flex-wrap gap-2 mt-2">
                {machine?.users?.map((user, index) => (
                  <span key={index} className="px-2 py-1 bg-blue-100 text-blue-700 rounded text-sm">
                    {user}
                  </span>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );

  const renderContainersView = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Running Containers</CardTitle>
          <CardDescription>
            Manage containers running on this machine
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr className="border-b">
                  <th className="text-left p-2 font-medium">Name</th>
                  <th className="text-left p-2 font-medium">Image</th>
                  <th className="text-left p-2 font-medium">Status</th>
                  <th className="text-left p-2 font-medium">Ports</th>
                  <th className="text-left p-2 font-medium">Actions</th>
                </tr>
              </thead>
              <tbody>
                {containers.map((container) => {
                  const isActionLoading = containerActions.some(action => action.id === container.id);
                  return (
                    <tr key={container.id} className="border-b hover:bg-gray-50">
                      <td className="p-2 font-medium">{container.name}</td>
                      <td className="p-2 text-gray-600">{container.image}</td>
                      <td className="p-2">
                        <span className={`px-2 py-1 rounded text-sm ${
                          container.status === "running"
                            ? "bg-green-100 text-green-700"
                            : "bg-red-100 text-red-700"
                        }`}>
                          {container.status}
                        </span>
                      </td>
                      <td className="p-2 text-gray-600">{container.ports}</td>
                      <td className="p-2">
                        <div className="flex gap-1">
                          {container.status === "stopped" ? (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleContainerAction(container.id, "start")}
                              disabled={isActionLoading}
                              className="h-8 w-8 p-0"
                            >
                              <Play className="h-4 w-4" />
                            </Button>
                          ) : (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleContainerAction(container.id, "stop")}
                              disabled={isActionLoading}
                              className="h-8 w-8 p-0"
                            >
                              <Square className="h-4 w-4" />
                            </Button>
                          )}
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleContainerAction(container.id, "restart")}
                            disabled={isActionLoading}
                            className="h-8 w-8 p-0"
                          >
                            <RotateCcw className="h-4 w-4" />
                          </Button>
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const renderMetricsView = () => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">CPU Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics?.cpu.current}%</div>
            <div className="text-xs text-muted-foreground">Avg: {metrics?.cpu.average}% | Peak: {metrics?.cpu.peak}%</div>
            <div className="w-full bg-gray-200 rounded-full h-2 mt-2">
              <div
                className="bg-blue-600 h-2 rounded-full"
                style={{ width: `${metrics?.cpu.current}%` }}
              ></div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Memory Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics?.memory.current}%</div>
            <div className="text-xs text-muted-foreground">Avg: {metrics?.memory.average}% | Peak: {metrics?.memory.peak}%</div>
            <div className="w-full bg-gray-200 rounded-full h-2 mt-2">
              <div
                className="bg-green-600 h-2 rounded-full"
                style={{ width: `${metrics?.memory.current}%` }}
              ></div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Network</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-lg font-bold">{metrics?.network.download} Mbps</div>
            <div className="text-xs text-muted-foreground">↓ Down | ↑ Up: {metrics?.network.upload} Mbps</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Uptime</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics?.uptime}%</div>
            <div className="text-xs text-muted-foreground">Last 30 days</div>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Resource Usage Over Time</CardTitle>
            <CardDescription>CPU and Memory usage trends</CardDescription>
          </CardHeader>
          <CardContent className="h-64 flex items-center justify-center bg-gray-50 rounded">
            <div className="text-gray-500">Chart placeholder - CPU & Memory trends</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Container Performance</CardTitle>
            <CardDescription>Individual container resource consumption</CardDescription>
          </CardHeader>
          <CardContent className="h-64 flex items-center justify-center bg-gray-50 rounded">
            <div className="text-gray-500">Chart placeholder - Container metrics</div>
          </CardContent>
        </Card>
      </div>
    </div>
  );

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Loading machine details...</div>
      </div>
    );
  }

  if (!machine) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-xl font-semibold mb-2">Machine not found</h2>
          <Button onClick={() => navigate("/dashboard")} className="flex items-center gap-2">
            <ArrowLeft className="h-4 w-4" />
            Back to Dashboard
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="flex">
        {/* Sidebar */}
        <div className="w-64 bg-white border-r border-gray-200 min-h-screen">
          <div className="p-6">
            <Button
              variant="ghost"
              onClick={() => navigate("/dashboard")}
              className="mb-4 flex items-center gap-2 w-full justify-start"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to Dashboard
            </Button>

            <h2 className="text-lg font-semibold mb-4">{machine.name}</h2>

            <nav className="space-y-2">
              <button
                onClick={() => setActiveView("machine")}
                className={`w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg transition-colors ${
                  activeView === "machine"
                    ? "bg-blue-100 text-blue-700 font-medium"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                <Server className="h-4 w-4" />
                Machine Details
              </button>

              <button
                onClick={() => setActiveView("containers")}
                className={`w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg transition-colors ${
                  activeView === "containers"
                    ? "bg-blue-100 text-blue-700 font-medium"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                <Box className="h-4 w-4" />
                Containers
              </button>

              <button
                onClick={() => setActiveView("metrics")}
                className={`w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg transition-colors ${
                  activeView === "metrics"
                    ? "bg-blue-100 text-blue-700 font-medium"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                <BarChart3 className="h-4 w-4" />
                Metrics
              </button>
            </nav>
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 p-8">
          <div className="max-w-6xl mx-auto">
            <div className="mb-6">
              <h1 className="text-2xl font-semibold capitalize mb-2">
                {activeView === "machine" && "Machine Details"}
                {activeView === "containers" && "Container Management"}
                {activeView === "metrics" && "Performance Metrics"}
              </h1>
              <p className="text-gray-600">
                {activeView === "machine" && "View and manage machine configuration and information"}
                {activeView === "containers" && "Monitor and control containers running on this machine"}
                {activeView === "metrics" && "Real-time performance monitoring and analytics"}
              </p>
            </div>

            {activeView === "machine" && renderMachineView()}
            {activeView === "containers" && renderContainersView()}
            {activeView === "metrics" && renderMetricsView()}
          </div>
        </div>
      </div>
    </div>
  );
};

export default MachinePage;
