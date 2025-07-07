// Simple auth utility for storing token and role in localStorage
export function setAuth(token: string, role: string) {
  if (typeof window !== "undefined") {
    localStorage.setItem("hospos_token", token);
    localStorage.setItem("hospos_role", role);
  }
}

export function getAuth() {
  if (typeof window !== "undefined") {
    return {
      token: localStorage.getItem("hospos_token"),
      role: localStorage.getItem("hospos_role"),
    };
  }
  return { token: null, role: null };
}

export function clearAuth() {
  if (typeof window !== "undefined") {
    localStorage.removeItem("hospos_token");
    localStorage.removeItem("hospos_role");
  }
}
