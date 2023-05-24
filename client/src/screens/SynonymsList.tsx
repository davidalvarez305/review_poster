import React, { useContext, useEffect, useMemo, useState } from "react";
import { Synonym } from "../types/general";
import Layout from "../layout/Layout";
import EditModal from "../components/EditModal";
import { Box, Button, Table, Tbody, Th, Thead, Tr } from "@chakra-ui/react";
import useLoginRequired from "../hooks/useLoginRequired";
import TableRow from "../components/TableRow";
import RequestErrorMessage from "../components/RequestErrorMessage";
import useSynonymsController from "../hooks/useSynonymsController";
import { UserContext } from "../context/UserContext";

export const SynonymsList: React.FC = () => {
  useLoginRequired();

  const {
    synonyms,
    isLoading,
    error,
    updateUserSynonymsByWord,
    deleteUserSynonymByWord,
    updateUserSynonymByWord,
    getUserSynonymsByWord,
  } = useSynonymsController();

  const { user } = useContext(UserContext);
  const [editModal, setEditModal] = useState(false);
  const [editingSynonym, setEditingSynonym] = useState<Synonym | null>(null);
  const [bulkModal, setBulkModal] = useState(false);

  const word = useMemo((): string | undefined => {
    return window.location.pathname.split("/word/")[1];
  }, []);

  useEffect(() => {
    if (word && user.id) getUserSynonymsByWord(word);
  }, [getUserSynonymsByWord, word, user.id]);

  const columns = ["id", "synonym", "action"];

  function handleSubmitBulk(values: { input: string }) {
    if (word) {
      updateUserSynonymsByWord({ ...values }, word);
    }
    setBulkModal(false);
  }

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
                        deleteUserSynonymByWord(synonyms[i].id!, word);
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
                updateUserSynonymByWord(
                  { ...editingSynonym, synonym: values.input },
                  editingSynonym.word?.name
                );
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
              setEditModal={setBulkModal}
              editModal={bulkModal}
              editingItem={synonyms.map((op) => op.synonym).join("\n")}
              isLoading={isLoading}
              handleSubmit={(values) => {
                handleSubmitBulk(values);
              }}
            />
          </Box>
        )}
      </Box>
    </Layout>
  );
};
