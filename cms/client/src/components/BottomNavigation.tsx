import { Box, Button } from "@chakra-ui/react";
import React from "react";
import { centeredDiv } from "../utils/centeredDiv";
import { useNavigate } from "react-router-dom";
import { Formik, Form } from "formik";
import SimpleInputField from "./SimpleInputField";
import { BsArrowReturnRight } from "react-icons/bs";

interface BottomNavigationProps {
  path: string;
  message: string;
}

export const BottomNavigation: React.FC<BottomNavigationProps> = ({ message, path }) => {
  const navigate = useNavigate();
  function handleClick(values: { destination: string }) {
    navigate("/" + path + "/" + values.destination);
  }
  return (
    <Formik initialValues={{ destination: "" }} onSubmit={handleClick}>
      <Form>
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
            w: '400px'
          }}
        >
          <SimpleInputField label={message} name={"destination"} />
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
        </Box>
      </Form>
    </Formik>
  );
};
