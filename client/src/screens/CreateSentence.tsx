import { Box, Button } from "@chakra-ui/react";
import React from "react";
import LargeInputBox from "../components/LargeInputBox";
import { centeredDiv } from "../utils/centeredDiv";
import { Formik, Form } from "formik";
import useFetch from "../hooks/useFetch";
import FormSelectComponent from "../components/FormSelectComponent";
import Layout from "../layout/Layout";
import { BottomNavigation } from "../components/BottomNavigation";
import useLoginRequired from "../hooks/useLoginRequired";
import RequestErrorMessage from "../components/RequestErrorMessage";
import useTemplatesController from "../hooks/useTemplatesController";
import useParagraphsController from "../hooks/useParagraphsController";
import useSentencesController from "../hooks/useSentencesController";

export const CreateSentence: React.FC = () => {
  const { isLoading, error } = useFetch();
  const { templates } = useTemplatesController();
  const { paragraphs, getParagraphs } = useParagraphsController();
  const { createSentences } = useSentencesController();
  useLoginRequired();

  return (
    <Layout>
      <Formik
        initialValues={{ paragraph: "", template: "", sentence: "" }}
        onSubmit={(values, actions) => {
          createSentences(values, paragraphs, templates);
          actions.resetForm({
            values: {
              paragraph: "",
              template: "",
              sentence: "",
            },
          });
          getParagraphs();
        }}
      >
        <Form>
          <Box sx={{ ...centeredDiv, gap: 2, height: "100%", my: 5 }}>
            <Box
              sx={{
                ...centeredDiv,
                width: "100%",
                height: "20%",
                flexDirection: "row",
              }}
            >
              <FormSelectComponent
                name={"paragraph"}
                options={paragraphs.map((p) => p.name)}
              />
              <FormSelectComponent
                name={"template"}
                options={[...new Set(templates.map((t) => t.name))]}
              />
            </Box>
            <Box>
              <LargeInputBox label="Sentence" name="sentence" />
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
      </Formik>
      <BottomNavigation message={"Enter A Paragraph"} path={"paragraph"} />
      <RequestErrorMessage {...error} />
    </Layout>
  );
};
