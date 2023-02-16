import { Box } from "@chakra-ui/layout";
import React from "react";
import { useNavigate } from "react-router";

interface Props {
  navigatePage: string;
  text: string;
}

const LoginOrRegister: React.FC<Props> = ({ navigatePage, text }) => {
  const navigate = useNavigate();
  return (
    <Box
      p={2}
      display={"flex"}
      justifyContent={"flex-end"}
      alignItems={"center"}
      onClick={() => navigate(`/${navigatePage}`)}
      cursor={"pointer"}
      textColor={"blue.600"}
      fontWeight={"semibold"}
    >
      {text}
    </Box>
  );
};

export default LoginOrRegister;
