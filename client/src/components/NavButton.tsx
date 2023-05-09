import { Button } from "@chakra-ui/react";
import React from "react";
import { useNavigate } from "react-router-dom";

interface NavButtonProps {
  buttonText: string;
  destination: string;
}

export const NavButton: React.FC<NavButtonProps> = ({
  buttonText,
  destination,
}) => {
  const navigate = useNavigate();
  return (
    <Button
      variant={"outline"}
      colorScheme={buttonText === "User" ? "red" : "cyan"}
      onClick={() => navigate("../" + destination)}
      size={"md"}
    >
      {buttonText}
    </Button>
  );
};
