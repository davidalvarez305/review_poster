import React from "react";
import { NavBar } from "../components/NavBar";
import Wrapper from "./Wrapper";

interface Props {
  children: React.ReactNode;
}

const Layout: React.FC<Props> = ({ children }) => {
  return (
    <Wrapper>
      <NavBar />
      {children}
    </Wrapper>
  );
};

export default Layout;
