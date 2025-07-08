// Simple auth utility for storing token, role, and username in localStorage
export function setAuth(token: string, role: string, username: string) {
  if (typeof window !== "undefined") {
    localStorage.setItem("hospos_token", token);
    localStorage.setItem("hospos_role", role);
    localStorage.setItem("hospos_username", username);
  }
}

export function getAuth() {
  if (typeof window !== "undefined") {
    return {
      token: localStorage.getItem("hospos_token"),
      role: localStorage.getItem("hospos_role"),
      username: localStorage.getItem("hospos_username"),
    };
  }
  return { token: null, role: null, username: null };
}

export function clearAuth() {
  if (typeof window !== "undefined") {
    localStorage.removeItem("hospos_token");
    localStorage.removeItem("hospos_role");
    localStorage.removeItem("hospos_username");
  }
}
