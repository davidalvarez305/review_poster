import { useCallback, useState } from "react";
import { LOGOUT_ROUTE } from "../constants";
import { User } from "../types/general";
import axios from "axios";

export const emptyUser = {
  id: 0,
  username: "",
  password: "",
  email: "",
  token: null,
  auth_header_string: "",
};

export default function useAuth() {
  const [user, setUser] = useState<User>(emptyUser);

  const SetUser = useCallback((user: User) => {
    setUser(user);
  }, []);

  function Logout() {
    axios
      .get(LOGOUT_ROUTE, {
        withCredentials: true,
      })
      .then((res) => {
        if (res.data.data === "Logged out!") {
          setUser(emptyUser);
        }
      })
      .catch(console.error);
  }

  return { SetUser, Logout, user };
}
