import { Box, Button, Input } from "@chakra-ui/react";
import React from "react";
import { centeredDiv } from "../utils/centeredDiv";
import { useNavigate } from "react-router-dom";
import { useFormik } from "formik";
import { BsArrowReturnRight } from "react-icons/bs";

interface BottomNavigationProps {
  path: string;
  message: string;
  dropdownComponent?: React.ReactNode;
}

export const BottomNavigation: React.FC<BottomNavigationProps> = ({
  message,
  path,
  dropdownComponent,
}) => {
  const navigate = useNavigate();
  function handleClick(values: { destination: string }) {
    navigate("/" + path + "/" + values.destination);
  }
  const formik = useFormik({
    initialValues: {
      destination: "",
    },
    onSubmit: handleClick,
  });
  return (
    <form onSubmit={formik.handleSubmit}>
      <Box
        sx={{
          ...centeredDiv,
          borderRadius: 25,
          borderColor: "teal",
          borderWidth: 0.5,
          height: 100,
          flexDirection: "row",
          justifyContent: "space-around",
          gap: 5,
          w: "700px",
        }}
      >
        <Input
          w={250}
          onChange={formik.handleChange}
          value={formik.values.destination}
          placeholder={message}
          name={"destination"}
          id={"destination"}
        />
        <Button
          variant={"outline"}
          colorScheme={"blue"}
          size={"md"}
          type={"submit"}
          loadingText={"Submitting"}
          width={"20%"}
          leftIcon={<BsArrowReturnRight />}
        >
          Go
        </Button>
        {dropdownComponent}
      </Box>
    </form>
  );
};
