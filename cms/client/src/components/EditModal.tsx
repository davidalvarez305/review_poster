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

interface Props {
  handleSubmit: (values: { input: string }) => void;
  setEditModal: React.Dispatch<React.SetStateAction<boolean>>;
  editModal: boolean;
  editingItem: string;
  isLoading: boolean;
  selectComponent?: React.ReactElement;
}

const EditModal: React.FC<Props> = ({
  handleSubmit,
  setEditModal,
  editModal,
  editingItem,
  isLoading,
  selectComponent,
}) => {
  const finalRef: React.RefObject<any> = useRef(null);
  return (
    <Formik initialValues={{ input: editingItem }} onSubmit={handleSubmit}>
      {({ submitForm }) => (
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
              </ModalFooter>
            </ModalContent>
          </Modal>
        </Form>
      )}
    </Formik>
  );
};

export default EditModal;
