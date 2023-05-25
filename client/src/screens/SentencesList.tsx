import React, { useCallback, useContext, useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import Layout from "../layout/Layout";
import {
  Box,
  Button,
  FormLabel,
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
import ReactSelect from "react-select";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import { UserContext } from "../context/UserContext";
import useSentencesController from "../hooks/useSentencesController";
import useParagraphsController from "../hooks/useParagraphsController";

export const SentencesList: React.FC = () => {
  useLoginRequired();
  const { makeRequest } = useFetch();
  const navigate = useNavigate();
  const { user } = useContext(UserContext);
  const {
    updateSentences,
    isLoading,
    error,
    deleteSentence,
    bulkUpdateSentences,
    updateSentence,
    getSentencesByParagraph,
    getSentences,
    sentences
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
      getSentencesByParagraph(paragraph);
    } else {
      getSentences();
    }
  }, [getSentences, getSentencesByParagraph, paragraph]);

  const headers = ["id", "sentence", "action"];

  function handleSubmitBulk(values: { input: string }) {
    if (selectedParagraph) {
      updateSentences({ ...values }, paragraphs[selectedParagraph].id!, paragraphs[selectedParagraph].name);
      navigate("/paragraph/" + paragraphs[selectedParagraph].name);
    } else {
      bulkUpdateSentences({ ...values });
    }

    setBulkModal(false);
  }

  const SelectChangeParagraph = useCallback(() => {
    return (
      <Box sx={{ width: 400, my: 2 }}>
        <FormLabel>Re-assign to a new paragraph (or leave bank)</FormLabel>
        <ReactSelect
          name={"select change paragraph"}
          placeholder={"select change paragraph"}
          value={{
            value: selectedParagraph ? selectedParagraph : "",
            label: selectedParagraph
              ? capitalizeFirstLetter(paragraphs[selectedParagraph].name)
              : "",
          }}
          onChange={(e) => {
            setSelectedParagraph(Number(e?.value));
          }}
          aria-label={"select change paragraph"}
          options={paragraphs.map((paragraph) => {
            return {
              value: paragraph.id!,
              label: capitalizeFirstLetter(paragraph.name),
            };
          })}
        />
      </Box>
    );
  }, [selectedParagraph, paragraphs]);

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
                    onClickDelete={() => deleteSentence(sentences[i].id!)}
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
               if (editingSentence) {
                  updateSentence({ ...editingSentence, sentence: values.input });

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
                selectComponent={SelectChangeParagraph()}
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
