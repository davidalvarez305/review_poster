import { Box } from "@chakra-ui/react";
import React from "react";
import { NavButton } from "./NavButton";

export const NavBar: React.FC = () => {
  const pages = [
    {
      buttonText: "Home",
      destination: "",
    },
    {
      buttonText: "Dictionary",
      destination: "dictionary",
    },
    {
      buttonText: "Create Sentence",
      destination: "/create-sentence",
    },
    {
      buttonText: "Create Paragraph",
      destination: "/create-paragraph",
    },
    {
      buttonText: "Generate",
      destination: "/generate",
    },
    {
      buttonText: "User",
      destination: "/user",
    },
  ];
  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "center",
        gap: 5,
      }}
    >
      {pages.map((page, index) => (
        <React.Fragment key={index}>
          <NavButton {...page} />
        </React.Fragment>
      ))}
    </Box>
  );
};
