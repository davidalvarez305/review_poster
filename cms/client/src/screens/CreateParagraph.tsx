import { Box, Button, useToast } from "@chakra-ui/react";
import React, { useContext, useEffect, useState } from "react";
import LargeInputBox from "../components/LargeInputBox";
import { centeredDiv } from "../utils/centeredDiv";
import { Formik, Form } from "formik";
import useFetch from "../hooks/useFetch";
import Layout from "../layout/Layout";
import { BottomNavigation } from "../components/BottomNavigation";
import { getId } from "../utils/getId";
import { Template } from "../types/general";
import FormSelectComponent from "../components/FormSelectComponent";
import useLoginRequired from "../hooks/useLoginRequired";
import { UserContext } from "../context/UserContext";
import { SaveButton } from "../components/SaveButton";
import RequestErrorMessage from "../components/RequestErrorMessage";
import { USER_ROUTE } from "../constants";

interface Props {}

export const CreateParagraph: React.FC<Props> = () => {
  const { user } = useContext(UserContext);
  const { isLoading, makeRequest, cancelToken, error } = useFetch();
  const [templates, setTemplates] = useState<Template[]>([]);
  const toast = useToast();
  useLoginRequired();

  function handleSubmit(values: { paragraphs: string; template: string }) {
    const template_id = getId(values.template, templates, "name");

    if (!template_id) {
      toast({
        title: "Heads Up!",
        description: "Save the template before submitting paragraphs.",
        status: "warning",
        isClosable: true,
        duration: 5000,
        variant: "left-accent",
      });
      return;
    }

    const paragraphBody = values.paragraphs.split("\n").map((i) => {
      return { name: i, order: null, template_id, user_id: user.id };
    });
  
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/paragraph`,
        method: "POST",
        data: paragraphBody,
      },
      () => {
        toast({
          title: "Success!",
          description: "Paragraphs have been submitted",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
      }
    );
  }

  function handleSaveTemplate(values: {
    paragraphs: string;
    template: string;
  }) {
    if (values.template === "") {
      return;
    }
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/template`,
        method: "POST",
        data: {
          name: values.template,
          user_id: user.id,
        },
      },
      (res) => {
        toast({
          title: "Success!",
          description: "Template has been submitted",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
        setTemplates(res.data.data);
      }
    );
  }

  useEffect(() => {
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/template`,
      },
      (res) => {
        setTemplates(res.data.data);
      }
    );
    return () => {
      cancelToken.cancel();
    };
  }, [makeRequest, cancelToken, user.id]);

  return (
    <Layout>
      <Formik
        initialValues={{ paragraphs: "", template: "" }}
        onSubmit={(values, actions) => {
          handleSubmit(values);
          actions.resetForm({
            values: {
              paragraphs: "",
              template: "",
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
                  options={templates.map((t) => t.name)}
                  name={"template"}
                />
                <SaveButton
                  aria-label={"save"}
                  onClick={() => {
                    handleSaveTemplate(values);
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
