import { Box, Button } from "@chakra-ui/react";
import React, { useEffect } from "react";
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
  const { paragraphs, getUserParagraphsByTemplate } = useParagraphsController();
  const { createUserParagraphSentencesByTemplate } = useSentencesController();
  const { templates, getUserTemplates } = useTemplatesController();
  useLoginRequired();

  useEffect(() => {
    getUserTemplates();
  }, [getUserTemplates]);

  return (
    <Layout>
      <Formik
        initialValues={{ paragraph: 0, template: 0, sentence: "" }}
        onSubmit={(values, actions) => {
          const paragraph = paragraphs.filter(t => t.id === values.paragraph)[0];
          createUserParagraphSentencesByTemplate({ input: values.sentence }, paragraph, paragraph?.template!);

          actions.resetForm({
            values: {
              paragraph: 0,
              template: 0,
              sentence: "",
            },
          });
          
          if (paragraph && paragraph.template) getUserParagraphsByTemplate(paragraph.template.name);
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
                options={paragraphs}
              />
              <FormSelectComponent
                name={"template"}
                options={templates}
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
      <BottomNavigation message={"Enter A Paragraph"} path={"template/ReviewPost/paragraph"} />
      <RequestErrorMessage {...error} />
    </Layout>
  );
};
