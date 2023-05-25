import { Box, Button } from "@chakra-ui/react";
import React, { useEffect } from "react";
import LargeInputBox from "../components/LargeInputBox";
import { centeredDiv } from "../utils/centeredDiv";
import { Formik, Form } from "formik";
import Layout from "../layout/Layout";
import { BottomNavigation } from "../components/BottomNavigation";
import FormSelectComponent from "../components/FormSelectComponent";
import useLoginRequired from "../hooks/useLoginRequired";
import { SaveButton } from "../components/SaveButton";
import RequestErrorMessage from "../components/RequestErrorMessage";
import useParagraphsController from "../hooks/useParagraphsController";
import useTemplatesController from "../hooks/useTemplatesController";

export const CreateParagraph: React.FC = () => {
  const { createParagraphs, isLoading, error } = useParagraphsController();
  const { templates, getUserTemplates } = useTemplatesController();
  useLoginRequired();

  useEffect(() => {
    getUserTemplates();
  }, [getUserTemplates]);

  return (
    <Layout>
      <Formik
        initialValues={{ paragraphs: "", template: 0 }}
        onSubmit={(values, actions) => {
          createParagraphs(values);
          actions.resetForm({
            values: {
              paragraphs: "",
              template: 0,
            },
          });
        }}
      >
        {({ values }) => (
          <Form>
            <Box sx={{ ...centeredDiv, gap: 2, height: "100%", my: 5 }}>
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  alignItems: "center",
                  flexDirection: "row",
                  gap: 5,
                }}
              >
                <FormSelectComponent
                  options={templates}
                  name={"template"}
                />
                <SaveButton
                  aria-label={"save"}
                  onClick={() => {
                    // createTemplates(values);
                    alert("Creating templates is disabled right now.")
                  }}
                  isLoading={isLoading}
                />
              </Box>
              <Box>
                <LargeInputBox
                  label="Paragraphs"
                  name="paragraphs"
                  height={"50vh"}
                  width={"25vw"}
                />
              </Box>
              <Button
                variant={"outline"}
                colorScheme={"teal"}
                size={"md"}
                type={"submit"}
                isLoading={isLoading}
                loadingText={"Submitting"}
              >
                Submit
              </Button>
            </Box>
          </Form>
        )}
      </Formik>
      <BottomNavigation message={"Enter A Template"} path={"template"} />
      <RequestErrorMessage {...error} />
    </Layout>
  );
};
