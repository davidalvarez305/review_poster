import React, { useContext, useEffect, useMemo, useState } from "react";
import Layout from "../layout/Layout";
import {
  Box,
  Button,
  Table,
  Tbody,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import useLoginRequired from "../hooks/useLoginRequired";
import EditModal from "../components/EditModal";
import TableRow from "../components/TableRow";
import useFetch from "../hooks/useFetch";
import { Sentence } from "../types/general";
import RequestErrorMessage from "../components/RequestErrorMessage";
import { UserContext } from "../context/UserContext";
import useSentencesController from "../hooks/useSentencesController";
import useParagraphsController from "../hooks/useParagraphsController";

export const SentencesList: React.FC = () => {
  useLoginRequired();
  const { makeRequest } = useFetch();
  const { user } = useContext(UserContext);
  const {
    sentences,
    isLoading,
    error,
    updateUserParagraphSentencesByTemplate,
    deleteUserParagraphSentenceByTemplate,
    updateParagraphSentenceByTemplate,
    getUserParagraphSentencesByTemplate,
  } = useSentencesController();
  const { paragraphs, getUserParagraphsByTemplate } = useParagraphsController();

  const [bulkModal, setBulkModal] = useState(false);
  const [editModal, setEditModal] = useState(false);
  const [editingSentence, setEditingSentence] = useState<Sentence | null>(null);
  const [selectedParagraph, setSelectedParagraph] = useState<number | null>(
    null
  );
  const paragraph = useMemo((): string | undefined => {
    return window.location.pathname.split("/paragraph/")[1]
  }, []);
  const template = useMemo((): string | undefined => {
    return window.location.pathname.split("/template/")[1]
  }, []);

  useEffect(() => {
    if (bulkModal) {
      // @TODO --> FIX THIS LOGIC
      getUserParagraphsByTemplate("ReviewPost");
    } else {
      setSelectedParagraph(null);
    }
  }, [bulkModal, makeRequest, user.id, getUserParagraphsByTemplate]);

  useEffect(() => {
    if (paragraph) {
      getUserParagraphSentencesByTemplate(paragraph, "ReviewPost");
    }
  }, [paragraph, getUserParagraphSentencesByTemplate]);

  const headers = ["id", "sentence", "action"];

  function handleSubmitBulk(values: { input: string }) {
    if (selectedParagraph) {
      updateUserParagraphSentencesByTemplate({ ...values }, paragraphs[selectedParagraph].template?.name!, paragraphs[selectedParagraph].name);
    }

    setBulkModal(false);
  }

  return (
    <Layout>
      <Box sx={{ my: 5 }}>
        <Button
          variant={"outline"}
          colorScheme={"green"}
          onClick={() => setBulkModal(!bulkModal)}
        >
          Bulk
        </Button>
      </Box>
      {sentences.length > 0 && (
        <Box
          sx={{
            width: "80%",
            height: "450px",
            overflow: "scroll",
          }}
        >
          <Table>
            <Thead>
              <Tr>
                {headers.map((h, idx) => (
                  <React.Fragment key={idx}>
                    <Th>{h}</Th>
                  </React.Fragment>
                ))}
              </Tr>
            </Thead>
            <Tbody>
              {sentences.map((_, i) => (
                <React.Fragment key={i}>
                  <TableRow
                    columns={headers}
                    index={i}
                    items={sentences}
                    onClickDelete={() => {
                      if (paragraph && template) deleteUserParagraphSentenceByTemplate(sentences[i].id!, paragraph, template);
                    }}
                    onClickEdit={() => {
                      setEditModal(true);
                      setEditingSentence(sentences[i]);
                    }}
                  />
                </React.Fragment>
              ))}
            </Tbody>
          </Table>
          <RequestErrorMessage {...error} />
          {editModal && (
            <EditModal
              editModal={editModal}
              setEditModal={setEditModal}
              handleSubmit={(values) => {
               if (editingSentence && template && paragraph) {
                updateParagraphSentenceByTemplate({ ...editingSentence, sentence: values.input }, template, paragraph);

                  setEditModal(false);
                  setEditingSentence(null);
               }
              }}
              editingItem={editingSentence?.sentence || ""}
              isLoading={isLoading}
            />
          )}
          {bulkModal && (
            <Box sx={{ my: 5 }}>
              <EditModal
                setEditModal={setBulkModal}
                editModal={bulkModal}
                editingItem={sentences.map((op) => op.sentence).join("\n")}
                isLoading={isLoading}
                handleSubmit={handleSubmitBulk}
              />
            </Box>
          )}
        </Box>
      )}
    </Layout>
  );
};
