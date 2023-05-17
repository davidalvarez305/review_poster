import React, { useCallback, useContext, useEffect, useState } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "../hooks/useFetch";
import { useLocation } from "react-router-dom";
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
import { Paragraph, Template } from "../types/general";
import RequestErrorMessage from "../components/RequestErrorMessage";
import { createUpdateParagraphsFactory } from "../utils/createUpdateParagraphsFactory";
import { UserContext } from "../context/UserContext";
import ReactSelect from "react-select";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import useParagraphsController from "../hooks/useParagraphsController";

interface ParagraphsListProps {}

export const ParagraphsList: React.FC<ParagraphsListProps> = () => {
  useLoginRequired();
  const { makeRequest } = useFetch();
  const { user } = useContext(UserContext);
  const {
    updateParagraphs,
    bulkUpdateParagraphs,
    deleteParagraph,
    paragraphs,
    isLoading,
    error,
  } = useParagraphsController();

  const [editModal, setEditModal] = useState(false);
  const [editingParagraph, setEditingParagraph] = useState("");
  const [bulkModal, setBulkModal] = useState(false);
  const [templates, setTemplates] = useState<Template[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<number | null>(null);

  useEffect(() => {
    if (bulkModal) {
      makeRequest(
        {
          url: USER_ROUTE + `/${user.id}/template`,
        },
        (res) => {
          setTemplates(res.data.data);
        }
      );
    }
    if (!bulkModal) {
      setSelectedTemplate(null);
    }
  }, [bulkModal, makeRequest, user.id]);

  const columns = ["id", "name", "order", "action"];

  function handleSubmit(values: { input: string }) {
    updateParagraphs(values);

    setEditModal(false);
    setEditOption(null);
    setEditingParagraph("");
  }

  function handleSubmitBulk(values: { input: string }) {
    const paragraphs = values.input.split("\n");

    let template_id = options[0].template_id;
    let templateString = template;
    let method = "POST";
    let body = paragraphs.map((name) => {
      return { name, template_id: options[0].template_id };
    });
    let route =
      USER_ROUTE + `/${user.id}/paragraph/bulk?template=${templateString}`;

    // Change request format if user selected a template.
    if (selectedTemplate) {
      template_id = templates[selectedTemplate].id;
      templateString = templates[selectedTemplate].name;
      body = createUpdateParagraphsFactory(
        options,
        paragraphs,
        template_id,
        templates[selectedTemplate]
      );
      route = USER_ROUTE + `/${user.id}/paragraph?template=${templateString}`;

      updateParagraphs({ ...values }, template_id, template: templateString);
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
          description: "Paragraphs have been successfully submitted.",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
        setBulkModal(false);
        setOptions(res.data.data);
      }
    );
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
          options={templates.map((op, index) => {
            return {
              value: index,
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
                    setEditingParagraph(paragraphs[i].name);
                    setEditOption(paragraphs[i]);
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
              updateParagraphs(values);

              setEditModal(false);
              setEditOption(null);
              setEditingParagraph("");
            }}
            editingItem={editingParagraph}
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
