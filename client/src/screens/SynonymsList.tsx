import React, { useCallback, useContext, useEffect, useState } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "../hooks/useFetch";
import { Synonym, Word } from "../types/general";
import { useLocation, useNavigate } from "react-router-dom";
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
  useToast,
} from "@chakra-ui/react";
import useLoginRequired from "../hooks/useLoginRequired";
import TableRow from "../components/TableRow";
import RequestErrorMessage from "../components/RequestErrorMessage";
import ReactSelect from "react-select";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import { createUpdateSynonyms } from "../utils/createUpdateSynonyms";
import { UserContext } from "../context/UserContext";

export const SynonymsList: React.FC = () => {
  useLoginRequired();
  const { user } = useContext(UserContext);
  const location = useLocation();
  const word = location.pathname.split("/word/")[1];
  const { makeRequest, isLoading, cancelToken, error } = useFetch();
  const toast = useToast();
  const navigate = useNavigate();

  const [options, setOptions] = useState<Synonym[]>([]);
  const [editModal, setEditModal] = useState(false);
  const [editOption, setEditOption] = useState<Synonym | null>();
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

  useEffect(() => {
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/synonym?word=${word}`,
      },
      (res) => {
        setOptions(res.data.data);
      }
    );
    return () => {
      cancelToken.cancel();
    };
  }, [editModal, word, makeRequest, cancelToken, user.id]);

  const columns = ["id", "synonym", "action"];

  function handleDelete(id: number) {
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/synonym?word=${word}&synonyms=${[id]}`,
        method: "DELETE",
      },
      (res) => {
        setOptions(res.data.data);
      }
    );
  }

  function handleSubmit(values: { input: string }) {
    const synonyms = values.input.split("\n");
    const body = synonyms.map((synonym) => {
      return { id: editOption?.id, word_id: editOption?.word_id, synonym };
    });
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/synonym/${editOption?.id}?word=${word}`,
        method: "PUT",
        data: body,
      },
      (res) => {
        toast({
          title: "Success!",
          description: "Synonym has been successfully submitted.",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
        setOptions(res.data.data);
        setEditOption(null);
        setEditingSynonym("");
      }
    );
    setEditModal(false);
  }

  function handleSubmitBulk(values: { input: string }) {
    const synonyms = values.input.split("\n");

    let word_id = options[0].word_id;
    let wordString = word;
    let method = "POST";
    let body = synonyms.map((synonym) => {
      return { synonym, word_id };
    });
    let route = USER_ROUTE + `/${user.id}/synonym/bulk?word=${wordString}`

    // Change request format if user selected a word.
    if (selectedWord) {
      word_id = words[selectedWord].id;
      wordString = words[selectedWord].name;
      method = "PUT";
      body = createUpdateSynonyms(options, synonyms, word_id, words[selectedWord]);
      route = USER_ROUTE + `/${user.id}/synonym?word=${wordString}`
    }

    makeRequest(
      {
        url: route,
        method: method,
        data: body,
      },
      (res) => {
        toast({
          title: "Success!",
          description: "Synonyms have been successfully submitted.",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
        setEditModal(false);
        setOptions(res.data.data);
        setBulkModal(false);
        if (selectedWord) {
          navigate("/word/" + wordString);
        }
      }
    );
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
              {options.map((_, i) => (
                <React.Fragment key={i}>
                  <TableRow
                    columns={columns}
                    index={i}
                    items={options}
                    onClickEdit={() => {
                      setEditModal(true);
                      setEditingSynonym(options[i].synonym);
                      setEditOption(options[i]);
                    }}
                    onClickDelete={() => handleDelete(options[i].id!)}
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
            handleSubmit={handleSubmit}
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
              editingItem={options.map((op) => op.synonym).join("\n")}
              isLoading={isLoading}
              handleSubmit={handleSubmitBulk}
            />
          </Box>
        )}
      </Box>
    </Layout>
  );
};
