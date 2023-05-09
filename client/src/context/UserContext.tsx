import React, { createContext } from "react";
import useAuth from "../hooks/useAuth";

export const UserContext = createContext<any>(null);

interface Props {
  children: React.ReactNode;
}

export const UserProvider = ({ children }: Props) => {
  const { SetUser, Logout, user } = useAuth();
  return (
    <UserContext.Provider value={{ SetUser, Logout, user }}>
      {children}
    </UserContext.Provider>
  );
};
