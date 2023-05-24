import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import useFetch from "../hooks/useFetch";
import { Synonym } from "../types/general";
import { useNavigate } from "react-router-dom";
import Layout from "../layout/Layout";
import EditModal from "../components/EditModal";
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
import TableRow from "../components/TableRow";
import RequestErrorMessage from "../components/RequestErrorMessage";
import ReactSelect from "react-select";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import useSynonymsController from "../hooks/useSynonymsController";
import { UserContext } from "../context/UserContext";
import useWordsController from "../hooks/useWordsController";

export const SynonymsList: React.FC = () => {
  useLoginRequired();

  const {
    updateSynonyms,
    synonyms,
    isLoading,
    error,
    bulkUpdateSynonyms,
    deleteSynonym,
    updateSynonym,
    getUserSynonymsByWord
  } = useSynonymsController();

  const { makeRequest } = useFetch();
  const navigate = useNavigate();
  const { user } = useContext(UserContext);
  const [editModal, setEditModal] = useState(false);
  const [editingSynonym, setEditingSynonym] = useState<Synonym | null>(null);
  const [bulkModal, setBulkModal] = useState(false);
  const [selectedWord, setSelectedWord] = useState<number | null>(null);
  const { words, getWords } = useWordsController();

  const word = useMemo((): string | undefined => {
    return window.location.pathname.split("/word/")[1];
  }, []);

  useEffect(() => {
    if (bulkModal) {
      getWords();
    }
    if (!bulkModal) {
      setSelectedWord(null);
    }
  }, [bulkModal, makeRequest, user.id, getWords]);

  useEffect(() => {
    if (word && user.id) getUserSynonymsByWord(word);
  }, [getUserSynonymsByWord, word, user.id]);

  const columns = ["id", "synonym", "action"];

  function handleSubmitBulk(values: { input: string }) {
    if (selectedWord) {
      updateSynonyms(
        { ...values },
        words[selectedWord!].id,
        words[selectedWord!].name
      );
      navigate("/word/" + word);
    } else if (word) {
      bulkUpdateSynonyms({ ...values }, word);
    }
    setBulkModal(false);
  }

  const SelectChangeWord = useCallback(() => {
    return (
      <Box sx={{ width: 400, my: 2 }}>
        <FormLabel>Select a new word (or leave bank)</FormLabel>
        <ReactSelect
          name={"select change word"}
          placeholder={"select change word"}
          value={{
            value: selectedWord ? selectedWord : "",
            label: selectedWord
              ? capitalizeFirstLetter(words[selectedWord].name)
              : "",
          }}
          onChange={(e) => {
            setSelectedWord(Number(e?.value));
          }}
          aria-label={"select change word"}
          options={words.map((op) => {
            return {
              value: op.id,
              label: capitalizeFirstLetter(op.name),
            };
          })}
        />
      </Box>
    );
  }, [selectedWord, words]);

  return (
    <Layout>
      <Box
        sx={{
          my: 5,
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <Button
          variant={"outline"}
          colorScheme={"green"}
          onClick={() => setBulkModal(true)}
        >
          Bulk
        </Button>
        <Box
          sx={{
            width: "1020px",
            height: "450px",
            overflow: "scroll",
          }}
        >
          <Table>
            <Thead>
              <Tr>
                {columns.map((column, idx) => (
                  <React.Fragment key={idx}>
                    <Th>{column}</Th>
                  </React.Fragment>
                ))}
              </Tr>
            </Thead>
            <Tbody>
              {synonyms.map((_, i) => (
                <React.Fragment key={i}>
                  <TableRow
                    columns={columns}
                    index={i}
                    items={synonyms}
                    onClickEdit={() => {
                      setEditModal(true);
                      setEditingSynonym(synonyms[i]);
                    }}
                    onClickDelete={() => {
                      if (word) {
                        deleteSynonym(synonyms[i].id!, word);
                      }
                    }}
                  />
                </React.Fragment>
              ))}
            </Tbody>
          </Table>
        </Box>
        <RequestErrorMessage {...error} />
        {editModal && (
          <EditModal
            editModal={editModal}
            setEditModal={setEditModal}
            handleSubmit={(values) => {
              if (editingSynonym && editingSynonym.word) {
                updateSynonym({ ...editingSynonym, synonym: values.input }, editingSynonym.word?.name);
                setEditingSynonym(null);
                setEditModal(false);
              }
            }}
            editingItem={editingSynonym?.synonym || ""}
            isLoading={isLoading}
          />
        )}
        {bulkModal && (
          <Box sx={{ my: 5 }}>
            <EditModal
              selectComponent={SelectChangeWord()}
              setEditModal={setBulkModal}
              editModal={bulkModal}
              editingItem={synonyms.map((op) => op.synonym).join("\n")}
              isLoading={isLoading}
              handleSubmit={(values) => {
                handleSubmitBulk(values);
              }}
              selectedWord={
                selectedWord ? words[selectedWord!].name : undefined
              }
            />
          </Box>
        )}
      </Box>
    </Layout>
  );
};
