import { Box, Button, useToast } from "@chakra-ui/react";
import React, { useContext, useEffect, useState } from "react";
import LargeInputBox from "../components/LargeInputBox";
import { centeredDiv } from "../utils/centeredDiv";
import { Formik, Form } from "formik";
import useFetch from "../hooks/useFetch";
import { SENTENCE_ROUTE } from "../constants";
import FormSelectComponent from "../components/FormSelectComponent";
import Layout from "../layout/Layout";
import { BottomNavigation } from "../components/BottomNavigation";
import { JoinedParagraph, Paragraph, Template } from "../types/general";
import { getId } from "../utils/getId";
import useLoginRequired from "../hooks/useLoginRequired";
import { UserContext } from "../context/UserContext";
import RequestErrorMessage from "../components/RequestErrorMessage";

interface Props {}

export const CreateSentence: React.FC<Props> = () => {
  const { user } = useContext(UserContext);
  const { isLoading, makeRequest, cancelToken, error } = useFetch();
  const [paragraphs, setParagraphs] = useState<Array<Paragraph>>([]);
  const [templates, setTemplates] = useState<Array<Template>>([]);
  const toast = useToast();
  useLoginRequired();

  function handleSubmit(values: {
    paragraph: string;
    template: string;
    sentence: string;
  }) {
    const paragraph_id = getId(values.paragraph, paragraphs, "name");
    const template_id = getId(values.template, templates, "name");
    const sentenceBody = values.sentence.split("\n").map((sentence) => {
      return {
        paragraph_id,
        template_id,
        sentence,
      };
    });

    if (!paragraph_id || !template_id) {
      toast({
        title: "Heads Up",
        description:
          "Before you can create sentences, you need to have a valid paragraph & template to submit to.",
        status: "warning",
        isClosable: true,
        duration: 5000,
        variant: "left-accent",
      });
      return;
    }
    makeRequest(
      {
        url: SENTENCE_ROUTE,
        method: "POST",
        data: sentenceBody,
      },
      () => {
        toast({
          title: "Success!",
          description: "Sentences have been submitted",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
      }
    );
  }

  useEffect(() => {
    makeRequest(
      {
        url: SENTENCE_ROUTE + "/create-sentence",
      },
      (res) => {
        const response: JoinedParagraph[] = res.data.data;
        setParagraphs(() =>
          response.map((r) => {
            return {
              user_id: user.id,
              template_id: r.template_id,
              name: r.paragraph_name,
              id: r.paragraph_id,
            };
          })
        );
        setTemplates(() =>
          response.map((r) => {
            return {
              name: r.template_name,
              id: r.template_id,
              user_id: user.id,
            };
          })
        );
      }
    );
    return () => {
      cancelToken.cancel();
    };
  }, [makeRequest, cancelToken, user]);

  return (
    <Layout>
      <Formik
        initialValues={{ paragraph: "", template: "", sentence: "" }}
        onSubmit={(values, actions) => {
          handleSubmit(values);
          actions.resetForm({
            values: {
              paragraph: "",
              template: "",
              sentence: "",
            },
          });
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
