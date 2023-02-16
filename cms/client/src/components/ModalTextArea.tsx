import {
  Box,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Textarea,
} from "@chakra-ui/react";
import { useField } from "formik";
import React, { InputHTMLAttributes } from "react";

type Props = InputHTMLAttributes<HTMLTextAreaElement> & {
  label: string;
  name: string;
};

const ModalTextArea: React.FC<Props> = ({ label, name, size: _, ...props }) => {
  const [field, meta] = useField(name);
  const centeredStyles = {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    width: "100%",
  };
  return (
    <React.Fragment>
      <Box sx={{ ...centeredStyles, flexDirection: "column" }}>
        <FormControl>
          <FormLabel sx={centeredStyles} htmlFor={field.name}>
            {label}
          </FormLabel>
          <Textarea {...props} {...field} id={field.name} minH="300px" />
          {meta.error && meta.touched && (
            <FormErrorMessage>{meta.error}</FormErrorMessage>
          )}
        </FormControl>
      </Box>
    </React.Fragment>
  );
};

export default ModalTextArea;
