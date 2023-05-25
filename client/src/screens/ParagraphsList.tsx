import React, { useContext, useEffect, useMemo, useState } from "react";
import Layout from "../layout/Layout";
import EditModal from "../components/EditModal";
import { Box, Button, Table, Tbody, Th, Thead, Tr } from "@chakra-ui/react";
import useLoginRequired from "../hooks/useLoginRequired";
import TableRow from "../components/TableRow";
import { Paragraph } from "../types/general";
import RequestErrorMessage from "../components/RequestErrorMessage";
import { UserContext } from "../context/UserContext";
import useParagraphsController from "../hooks/useParagraphsController";

export const ParagraphsList: React.FC = () => {
  useLoginRequired();
  const { user } = useContext(UserContext);
  const {
    updateUserParagraphsByTemplate,
    updateUserParagraphByTemplate,
    getUserParagraphsByTemplate,
    deleteUserParagraphByWord,
    paragraphs,
    isLoading,
    error,
  } = useParagraphsController();
  const [editModal, setEditModal] = useState(false);
  const [editingParagraph, setEditingParagraph] = useState<Paragraph | null>(
    null
  );
  const [bulkModal, setBulkModal] = useState(false);
  const template = useMemo((): string | undefined => {
    return window.location.pathname.split("/template/")[1];
  }, []);

  useEffect(() => {
    if (template && user.id) getUserParagraphsByTemplate(template);
  }, [getUserParagraphsByTemplate, template, user.id]);

  const columns = ["id", "name", "order", "action"];

  function handleSubmitBulk(values: { input: string }) {
    if (template) {
      updateUserParagraphsByTemplate({ ...values }, template);
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
              {columns.map((column, idx) => (
                <React.Fragment key={idx}>
                  <Th>{column}</Th>
                </React.Fragment>
              ))}
            </Tr>
          </Thead>
          <Tbody>
            {paragraphs.map((_, i) => (
              <React.Fragment key={i}>
                <TableRow
                  columns={columns}
                  index={i}
                  items={paragraphs}
                  onClickEdit={() => {
                    setEditModal(true);
                    setEditingParagraph(paragraphs[i]);
                  }}
                  onClickDelete={() => {
                    if (template) deleteUserParagraphByWord(paragraphs[i].id!, template);
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
              if (editingParagraph && editingParagraph.template) {
                updateUserParagraphByTemplate(
                  { ...editingParagraph, name: values.input, template: null },
                  editingParagraph.template?.name
                );
                setEditingParagraph(null);
                setEditModal(false);
              }
            }}
            editingItem={editingParagraph?.name || ""}
            isLoading={isLoading}
          />
        )}
        {bulkModal && (
          <Box sx={{ my: 5 }}>
            <EditModal
              setEditModal={setBulkModal}
              editModal={bulkModal}
              editingItem={paragraphs.map((op) => op.name).join("\n")}
              isLoading={isLoading}
              handleSubmit={handleSubmitBulk}
            />
          </Box>
        )}
      </Box>
    </Layout>
  );
};
