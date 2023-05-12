import { useCallback, useState } from "react";
import { LOGOUT_ROUTE } from "../constants";
import { User } from "../types/general";
import useFetch from "./useFetch";
export const emptyUser = {
  id: 0,
  username: "",
  password: "",
  email: "",
  token: null
};

export default function useAuth() {
  const { makeRequest } = useFetch();
  const [user, setUser] = useState<User>(emptyUser);

  const SetUser = useCallback((user: User) => {
    setUser(user);
  }, []);

  function Logout() {
    makeRequest(
      {
        url: `${LOGOUT_ROUTE}`,
        method: "POST",
      },
      (res) => {
        if (res.data.data === "Logged out!") {
          setUser(emptyUser);
        }
      }
    );
  }

  return { SetUser, Logout, user };
}
