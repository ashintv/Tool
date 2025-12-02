import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { X } from "lucide-react";
import axios from "axios";

interface RegisterMachineRequest {
  name: string;
  users: number[];
  ip: string;
}

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
}

interface AddMachineModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

// Mock user data - replace with API call in the future
const mockUsers: User[] = [
  { id: 1, name: "John Doe", email: "john.doe@example.com", role: "Admin" },
  { id: 2, name: "Jane Smith", email: "jane.smith@example.com", role: "Developer" },
  { id: 3, name: "Mike Johnson", email: "mike.johnson@example.com", role: "DevOps" },
  { id: 4, name: "Sarah Wilson", email: "sarah.wilson@example.com", role: "Developer" },
  { id: 5, name: "David Brown", email: "david.brown@example.com", role: "Manager" },
  { id: 6, name: "Lisa Davis", email: "lisa.davis@example.com", role: "QA Engineer" },
];

export const AddMachineModal: React.FC<AddMachineModalProps> = ({ isOpen, onClose, onSuccess }) => {
  const [formData, setFormData] = useState<RegisterMachineRequest>({
    name: "",
    users: [],
    ip: "",
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [selectedUsers, setSelectedUsers] = useState<number[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  if (!isOpen) return null;

  const handleInputChange = (field: keyof RegisterMachineRequest, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleUserSelect = (userId: number) => {
    if (selectedUsers.includes(userId)) {
      const newSelectedUsers = selectedUsers.filter(id => id !== userId);
      setSelectedUsers(newSelectedUsers);
      setFormData(prev => ({
        ...prev,
        users: newSelectedUsers,
      }));
    } else {
      const newSelectedUsers = [...selectedUsers, userId];
      setSelectedUsers(newSelectedUsers);
      setFormData(prev => ({
        ...prev,
        users: newSelectedUsers,
      }));
    }
  };

  const getSelectedUserNames = () => {
    return selectedUsers
      .map(userId => mockUsers.find(user => user.id === userId)?.name)
      .filter(Boolean)
      .join(", ");
  };

  const filteredUsers = mockUsers.filter(user =>
    user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.email.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const removeUser = (userId: number) => {
    const newSelectedUsers = selectedUsers.filter(id => id !== userId);
    setSelectedUsers(newSelectedUsers);
    setFormData(prev => ({
      ...prev,
      users: newSelectedUsers,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name || !formData.ip || formData.users.length === 0) {
      alert("Please fill in all required fields");
      return;
    }

    setIsSubmitting(true);
    try {
      const response = await axios.post(
        "http://localhost:8080/api/user/machine/register",
        formData,
        {
          headers: {
            Authorization: localStorage.getItem("token"),
            "Content-Type": "application/json",
          },
        }
      );

      if (response.status === 200) {
        onSuccess();
        onClose();
        // Reset form
        setFormData({ name: "", users: [], ip: "" });
        setSelectedUsers([]);
      }
    } catch (error) {
      console.error("Error registering machine:", error);
      alert("Failed to register machine. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setFormData({ name: "", users: [], ip: "" });
    setSelectedUsers([]);
    setSearchTerm("");
    setIsDropdownOpen(false);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50"
        onClick={() => {
          setIsDropdownOpen(false);
          handleClose();
        }}
      />

      {/* Modal */}
      <div className="relative z-50 w-full max-w-md mx-4 bg-background rounded-lg border shadow-lg">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-lg font-semibold">Add New Machine</h2>
          <Button
            variant="ghost"
            size="icon"
            onClick={handleClose}
            className="h-6 w-6"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Machine Name *</Label>
            <Input
              id="name"
              type="text"
              placeholder="Enter machine name"
              value={formData.name}
              onChange={(e) => handleInputChange("name", e.target.value)}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="ip">IP Address *</Label>
            <Input
              id="ip"
              type="text"
              placeholder="192.168.1.100"
              value={formData.ip}
              onChange={(e) => handleInputChange("ip", e.target.value)}
              required
            />
          </div>

          <div className="space-y-2">
            <Label>Assign Users *</Label>

            {/* Selected Users Display */}
            {selectedUsers.length > 0 && (
              <div className="flex flex-wrap gap-2 p-2 border rounded-md bg-muted/50">
                {selectedUsers.map((userId) => {
                  const user = mockUsers.find(u => u.id === userId);
                  return user ? (
                    <div key={userId} className="flex items-center gap-1 bg-primary text-primary-foreground px-2 py-1 rounded-md text-sm">
                      <span>{user.name}</span>
                      <button
                        type="button"
                        onClick={() => removeUser(userId)}
                        className="ml-1 hover:bg-primary-foreground/20 rounded-full p-0.5"
                      >
                        <X className="h-3 w-3" />
                      </button>
                    </div>
                  ) : null;
                })}
              </div>
            )}

            {/* Multi-Select Dropdown */}
            <div className="relative">
              <div className="relative">
                <Input
                  type="text"
                  placeholder="Search and select users..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  onFocus={() => setIsDropdownOpen(true)}
                  className="pr-8"
                />
                <button
                  type="button"
                  onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                  className="absolute right-2 top-1/2 transform -translate-y-1/2"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
              </div>

              {isDropdownOpen && (
                <div className="absolute z-10 w-full mt-1 bg-background border rounded-md shadow-lg max-h-60 overflow-y-auto">
                  {filteredUsers.length === 0 ? (
                    <div className="p-3 text-sm text-muted-foreground text-center">
                      No users found
                    </div>
                  ) : (
                    filteredUsers.map((user) => (
                      <div
                        key={user.id}
                        className="flex items-center space-x-3 p-3 hover:bg-muted cursor-pointer border-b last:border-b-0"
                        onClick={() => {
                          handleUserSelect(user.id);
                          setSearchTerm("");
                        }}
                      >
                        <input
                          type="checkbox"
                          checked={selectedUsers.includes(user.id)}
                          onChange={() => {}}
                          className="rounded border-gray-300"
                        />
                        <div className="flex-1">
                          <div className="font-medium text-sm">{user.name}</div>
                          <div className="text-xs text-muted-foreground">
                            {user.email} • {user.role}
                          </div>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              )}
            </div>

            <p className="text-xs text-muted-foreground">
              Search by name or email. {selectedUsers.length} user(s) selected.
            </p>
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-3 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Adding..." : "Add Machine"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
