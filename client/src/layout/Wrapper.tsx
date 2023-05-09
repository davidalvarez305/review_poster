import { Box } from "@chakra-ui/layout";
import React from "react";
import { centeredDiv } from "../utils/centeredDiv";

interface Props {
  children: React.ReactNode;
}

const Wrapper: React.FC<Props> = ({ children }) => {
  return <Box sx={{ ...centeredDiv, my: 10 }}>{children}</Box>;
};

export default Wrapper;
