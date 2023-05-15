import {
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Button,
  Modal,
} from "@chakra-ui/react";
import { Form, Formik } from "formik";
import React, { useRef } from "react";
import ModalTextArea from "./ModalTextArea";
import { API_ROUTE } from "../constants";
import { createTag } from "../utils/createTag";
import useFetch from "../hooks/useFetch";

interface Props {
  handleSubmit: (values: { input: string }) => void;
  setEditModal: React.Dispatch<React.SetStateAction<boolean>>;
  editModal: boolean;
  editingItem: string;
  isLoading: boolean;
  selectComponent?: React.ReactElement;
  selectedWord?: string;
}

const EditModal: React.FC<Props> = ({
  handleSubmit,
  setEditModal,
  editModal,
  editingItem,
  isLoading,
  selectComponent,
  selectedWord
}) => {
  const { makeRequest, isLoading: isPulling } = useFetch();
  const finalRef: React.RefObject<any> = useRef(null);

  function handlePullFromChatGPT(
    setFieldValue: (
      field: string,
      value: any,
      shouldValidate?: boolean | undefined
    ) => void
  ) {
    const tagWord = selectedWord ? selectedWord : window.location.pathname.split("/word/")[1];
    const tag = createTag(tagWord);
    makeRequest(
      {
        url: API_ROUTE + `/ai/tags?tag=${encodeURIComponent(tag)}`,
      },
      async (res) => {
        setFieldValue("input", res.data.data.join("\n"));
      }
    );
  }

  return (
    <Formik initialValues={{ input: editingItem }} onSubmit={handleSubmit}>
      {({ submitForm, setFieldValue }) => (
        <Form>
          <Modal
            size="3xl"
            finalFocusRef={finalRef}
            isOpen={editModal}
            onClose={() => setEditModal(false)}
          >
            <ModalOverlay />
            <ModalContent>
              <ModalHeader>{`Editing...`}</ModalHeader>
              <ModalCloseButton />
              <ModalBody>
                {selectComponent}
                <ModalTextArea label={"Create Edits"} name={"input"} />
              </ModalBody>
              <ModalFooter>
                <Button
                  colorScheme="blue"
                  mr={3}
                  loadingText="Submitting"
                  isLoading={isLoading}
                  type="submit"
                  onClick={submitForm}
                >
                  Submit
                </Button>
                <Button
                  colorScheme="red"
                  mr={3}
                  onClick={() => setEditModal(false)}
                >
                  Close
                </Button>
                {selectComponent && <Button
                  variant={"outline"}
                  colorScheme={"red"}
                  size={"md"}
                  type={"button"}
                  isLoading={isPulling}
                  loadingText={"Pulling"}
                  onClick={() =>
                    handlePullFromChatGPT(setFieldValue)
                  }
                >
                  Pull From ChatGPT
                </Button>}
              </ModalFooter>
            </ModalContent>
          </Modal>
        </Form>
      )}
    </Formik>
  );
};

export default EditModal;
