import React, { useCallback, useContext, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
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
  useToast,
} from "@chakra-ui/react";
import useLoginRequired from "../hooks/useLoginRequired";
import EditModal from "../components/EditModal";
import TableRow from "../components/TableRow";
import { PARAGRAPH_ROUTE, SENTENCE_ROUTE } from "../constants";
import useFetch from "../hooks/useFetch";
import { Paragraph, Sentence } from "../types/general";
import RequestErrorMessage from "../components/RequestErrorMessage";
import ReactSelect from "react-select";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import { createUpdateSentences } from "../utils/createUpdateSentences";
import { UserContext } from "../context/UserContext";

export const SentencesList: React.FC = () => {
  useLoginRequired();
  const location = useLocation();
  const paragraph = location.pathname.split("/paragraph/")[1];
  const { makeRequest, isLoading, cancelToken, error } = useFetch();
  const navigate = useNavigate();
  const toast = useToast();
  const { user } = useContext(UserContext);

  const [bulkModal, setBulkModal] = useState(false);
  const [options, setOptions] = useState<Sentence[]>([]);
  const [editModal, setEditModal] = useState(false);
  const [editOption, setEditOption] = useState<Sentence | null>();
  const [editingSentence, setEditingSentence] = useState("");
  const [paragraphs, setParagraphs] = useState<Paragraph[]>([]);
  const [selectedParagraph, setSelectedParagraph] = useState<number | null>(
    null
  );

  useEffect(() => {
    if (bulkModal) {
      makeRequest(
        {
          url: PARAGRAPH_ROUTE,
        },
        (res) => {
          setParagraphs(res.data.data);
        }
      );
    }
    if (!bulkModal) {
      setSelectedParagraph(null);
    }
  }, [bulkModal, makeRequest]);

  useEffect(() => {
    makeRequest(
      {
        url: SENTENCE_ROUTE + `/${paragraph}`,
      },
      (res) => {
        setOptions(res.data.data);
      }
    );
  }, [editModal, paragraph, makeRequest, cancelToken]);

  const headers = ["id", "sentence", "action"];

  function handleDelete(id: number) {
    makeRequest(
      {
        url: SENTENCE_ROUTE + `/?sentences=${[id]}&paragraph=${paragraph}`,
        method: "DELETE",
      },
      (res) => {
        setOptions(res.data.data);
      }
    );
  }

  function handleSubmit(values: { input: string }) {
    const sentence = values.input.split("\n")[0];
    makeRequest(
      {
        url: SENTENCE_ROUTE,
        method: "PUT",
        data: {
          id: editOption?.id,
          paragraph_id: editOption?.paragraph_id,
          template_id: editOption?.template_id,
          sentence,
        },
      },
      (res) => {
        toast({
          title: "Success!",
          description: "Sentence has been successfully submitted.",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
        setEditModal(false);
        setOptions(res.data.data);
        setEditOption(null);
        setEditingSentence("");
      }
    );
  }

  function handleSubmitBulk(values: { input: string }) {
    const sentences = values.input.split("\n");

    let template_id = options[0].template_id;
    let paragraph_id = options[0].paragraph_id;
    let paragraphString = paragraph;
    let method = "POST";
    let body = sentences.map((sentence) => {
      return {
        paragraph_id,
        template_id,
        sentence,
        user_id: user.id,
      };
    });
    let route = SENTENCE_ROUTE + "/bulk/?paragraph=" + paragraphString;

    // Change request format if user selected a paragraph.
    if (selectedParagraph) {
      paragraph_id = paragraphs[selectedParagraph].id!;
      template_id = paragraphs[selectedParagraph].template_id;
      paragraphString = paragraphs[selectedParagraph].name;
      method = "PUT";
      body = createUpdateSentences(
        options,
        sentences,
        template_id,
        paragraph_id,
        user.id
      );
      route = SENTENCE_ROUTE + "?paragraph=" + paragraphString;
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
          description: "Sentences have been successfully submitted.",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
        setBulkModal(false);
        setOptions(res.data.data);
        if (selectedParagraph) {
          navigate("/paragraph/" + paragraphString);
        }
      }
    );
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
          options={paragraphs.map((paragraph, index) => {
            return {
              value: index,
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
      {options.length > 0 && (
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
              {options.map((_, i) => (
                <React.Fragment key={i}>
                  <TableRow
                    columns={headers}
                    index={i}
                    items={options}
                    onClickDelete={() => handleDelete(options[i].id!)}
                    onClickEdit={() => {
                      setEditModal(true);
                      setEditingSentence(options[i].sentence);
                      setEditOption(options[i]);
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
              handleSubmit={handleSubmit}
              editingItem={editingSentence}
              isLoading={isLoading}
            />
          )}
          {bulkModal && (
            <Box sx={{ my: 5 }}>
              <EditModal
                selectComponent={SelectChangeParagraph()}
                setEditModal={setBulkModal}
                editModal={bulkModal}
                editingItem={options.map((op) => op.sentence).join("\n")}
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
