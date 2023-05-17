import React, { useCallback, useContext, useEffect, useState } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "../hooks/useFetch";
import { Word } from "../types/general";
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

export const SynonymsList: React.FC = () => {
  useLoginRequired();

  const {
    updateSynonyms,
    synonyms,
    isLoading,
    error,
    setEditSingleSynonym,
    bulkUpdateSynonyms,
    deleteSynonym,
    word,
  } = useSynonymsController();

  const { makeRequest } = useFetch();
  const navigate = useNavigate();
  const { user } = useContext(UserContext);
  const [editModal, setEditModal] = useState(false);
  const [editingSynonym, setEditingSynonym] = useState("");
  const [bulkModal, setBulkModal] = useState(false);
  const [words, setWords] = useState<Word[]>([]);
  const [selectedWord, setSelectedWord] = useState<number | null>(null);

  useEffect(() => {
    if (bulkModal) {
      makeRequest(
        {
          url: USER_ROUTE + `/${user.id}/word`,
        },
        (res) => {
          setWords(res.data.data);
        }
      );
    }
    if (!bulkModal) {
      setSelectedWord(null);
    }
  }, [bulkModal, makeRequest, user.id]);

  const columns = ["id", "synonym", "action"];

  function handleSubmitBulk(values: { input: string }) {
    if (selectedWord) {
      updateSynonyms(
        { ...values },
        words[selectedWord!].id,
        words[selectedWord!].name
      );
    } else {
      bulkUpdateSynonyms({ ...values });
    }

    setEditModal(false);
    setBulkModal(false);
    if (selectedWord) {
      navigate("/word/" + word);
    }
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
          options={words.map((op, index) => {
            return {
              value: index,
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
                      setEditingSynonym(synonyms[i].synonym);
                      setEditSingleSynonym(synonyms[i]);
                    }}
                    onClickDelete={() => deleteSynonym(synonyms[i].id!)}
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
              if (word) {
                updateSynonyms(values, synonyms[0].id!, word);
                setEditingSynonym("");
                setEditModal(false);
              }
            }}
            editingItem={editingSynonym}
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
