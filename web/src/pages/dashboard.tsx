import React, { useEffect, useState } from "react";
import type { Machine } from "@/types/types";
import { MachineCard } from "@/components/machineCard";
import { AddMachineModal } from "@/components/AddMachineModal";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Plus } from "lucide-react";
import axios from "axios";

export const DashboardPage = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [machines, setMachines] = useState<Machine[]>([]);
  const [machineType, setMachineType] = useState<"owned" | "usable">("usable");
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);

  useEffect(() => {
    async function FetchMachines() {
      try {
        const endpoint =
          machineType === "owned"
            ? "http://localhost:8080/api/user/machine/owned"
            : "http://localhost:8080/api/user/machine/usable";

        const res = await axios.get(endpoint, {
          headers: {
            Authorization: localStorage.getItem("token"),
          },
        });
        console.log(res.data);
        const machinesData: Machine[] = [];

        res.data.machines.forEach((machine: any) => {
          machinesData.push({
            id: machine.ID,
            name: machine.Name,
            status: "online",
            IP: machine.IP,
          });
        });
        setMachines(machinesData);
      } catch (error) {
        console.error("There was an error fetching the machines!", error);
      }
    }
    FetchMachines();
  }, [machineType]);

  const handleConnect = (machine: Machine) => {
    console.log(`Connecting to machine: ${machine.name}`);
    setIsLoading(true);
    // Simulate API call
    setTimeout(() => {
      setIsLoading(false);
    }, 2000);
  };

  const handleDisconnect = (machine: Machine) => {
    console.log(`Disconnecting from machine: ${machine.name}`);
    setIsLoading(true);
    // Simulate API call
    setTimeout(() => {
      setIsLoading(false);
    }, 2000);
  };

  const handleEdit = (machine: Machine) => {
    console.log(`Editing machine: ${machine.name}`);
  };

  const handleDelete = (machine: Machine) => {
    console.log(`Deleting machine: ${machine.name}`);
  };

  const handleAddMachine = () => {
    setIsAddModalOpen(true);
  };

  const handleAddMachineSuccess = () => {
    // Refresh the machines list after successful addition
    const endpoint =
      machineType === "owned"
        ? "http://localhost:8080/api/user/machine/owned"
        : "http://localhost:8080/api/user/machine/usable";

    axios.get(endpoint, {
      headers: {
        Authorization: localStorage.getItem("token"),
      },
    })
    .then(res => {
      const machinesData: Machine[] = [];
      res.data.machines.forEach((machine: any) => {
        machinesData.push({
          id: machine.ID,
          name: machine.Name,
          status: "online",
          IP: machine.IP,
        });
      });
      setMachines(machinesData);
    })
    .catch(error => {
      console.error("Error refreshing machines:", error);
    });
  };

  return (
    <div className="dashboard-container">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8 flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-semibold text-foreground mb-2">Machine Dashboard</h1>
            <p className="text-muted-foreground">Manage and monitor your machines</p>
          </div>
          <div className="flex items-center gap-4">
            <Select value={machineType} onValueChange={(value: "owned" | "usable") => setMachineType(value)}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="owned">Owned</SelectItem>
                <SelectItem value="usable">Usable</SelectItem>
              </SelectContent>
            </Select>
            <Button onClick={handleAddMachine} className="flex items-center gap-2">
              <Plus className="size-4" />
              Add Machine
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {machines.map((machine) => (
            <MachineCard
              key={machine.id}
              machine={machine}
              onConnect={handleConnect}
              onDisconnect={handleDisconnect}
              onEdit={handleEdit}
              onDelete={handleDelete}
              isLoading={isLoading}
            />
          ))}
        </div>

        {machines.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No machines found. Add a machine to get started.</p>
          </div>
        )}
      </div>

      {/* Add Machine Modal */}
      <AddMachineModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onSuccess={handleAddMachineSuccess}
      />
    </div>
  );
};
