import React, { useCallback, useContext, useEffect, useState } from "react";
import useFetch from "../hooks/useFetch";
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
import { Paragraph } from "../types/general";
import RequestErrorMessage from "../components/RequestErrorMessage";
import { UserContext } from "../context/UserContext";
import ReactSelect from "react-select";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import useParagraphsController from "../hooks/useParagraphsController";
import useTemplatesController from "../hooks/useTemplatesController";

export const ParagraphsList: React.FC = () => {
  useLoginRequired();
  const { makeRequest } = useFetch();
  const { user } = useContext(UserContext);
  const {
    updateParagraphs,
    bulkUpdateParagraphs,
    deleteParagraph,
    updateParagraph,
    paragraphs,
    isLoading,
    error,
  } = useParagraphsController();
  const { templates, getTemplates } = useTemplatesController();

  const [editModal, setEditModal] = useState(false);
  const [editingParagraph, setEditingParagraph] = useState<Paragraph | null>(
    null
  );
  const [bulkModal, setBulkModal] = useState(false);
  const [selectedTemplate, setSelectedTemplate] = useState<number | null>(null);

  useEffect(() => {
    if (bulkModal) {
      getTemplates();
    }
    if (!bulkModal) {
      setSelectedTemplate(null);
    }
  }, [bulkModal, makeRequest, user.id, getTemplates]);

  const columns = ["id", "name", "order", "action"];

  function handleSubmitBulk(values: { input: string }) {
    if (selectedTemplate) {
      updateParagraphs(
        { ...values },
        templates[selectedTemplate].id,
        templates[selectedTemplate].name
      );
    } else {
      bulkUpdateParagraphs({ ...values });
    }
    setBulkModal(false);
  }

  const SelectChangeTemplate = useCallback(() => {
    return (
      <Box sx={{ width: 400, my: 2 }}>
        <FormLabel>Re-assign to a new template (or leave bank)</FormLabel>
        <ReactSelect
          name={"select change template"}
          placeholder={"select change template"}
          value={{
            value: selectedTemplate ? selectedTemplate : "",
            label: selectedTemplate
              ? capitalizeFirstLetter(templates[selectedTemplate].name)
              : "",
          }}
          onChange={(e) => {
            setSelectedTemplate(Number(e?.value));
          }}
          aria-label={"select change template"}
          options={templates.map((op) => {
            return {
              value: op.id,
              label: capitalizeFirstLetter(op.name),
            };
          })}
        />
      </Box>
    );
  }, [selectedTemplate, templates]);

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
                  onClickDelete={() => deleteParagraph(paragraphs[i].id!)}
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
              if (editingParagraph) {
                updateParagraph({
                  ...editingParagraph,
                  name: values.input,
                  template: null,
                });

                setEditModal(false);
                setEditingParagraph(null);
              }
            }}
            editingItem={editingParagraph?.name || ""}
            isLoading={isLoading}
          />
        )}
        {bulkModal && (
          <Box sx={{ my: 5 }}>
            <EditModal
              selectComponent={SelectChangeTemplate()}
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
