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
import { createUpdateParagraphs } from "../utils/createUpdateParagraphs";
import { UserContext } from "../context/UserContext";
import ReactSelect from "react-select";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";

interface ParagraphsListProps {}

export const ParagraphsList: React.FC<ParagraphsListProps> = () => {
  useLoginRequired();
  const location = useLocation();
  const template = location.pathname.split("/template/")[1];
  const { makeRequest, isLoading, cancelToken, error } = useFetch();
  const toast = useToast();
  const { user } = useContext(UserContext);

  const [options, setOptions] = useState<Paragraph[]>([]);
  const [editModal, setEditModal] = useState(false);
  const [editOption, setEditOption] = useState<Paragraph | null>();
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

  useEffect(() => {
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/paragraph?template=${template}`,
      },
      (res) => {
        setOptions(res.data.data);
      }
    );
    return () => {
      cancelToken.cancel();
    };
  }, [editModal, template, makeRequest, cancelToken, user.id]);
  const columns = ["id", "name", "order", "action"];

  function handleDelete(id: number) {
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/paragraph?paragraphs=${[id]}&template=${template}`,
        method: "DELETE",
      },
      (res) => {
        setOptions(res.data.data);
      }
    );
  }

  function handleSubmit(values: { input: string }) {
    const paragraph = values.input.split("\n")[0];
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/paragraph/`,
        method: "PUT",
        data: {
          id: editOption?.id,
          order: editOption?.order,
          name: paragraph,
          template_id: editOption?.template_id,
        },
      },
      (res) => {
        toast({
          title: "Success!",
          description: "Paragraph has been successfully submitted.",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
        setEditModal(false);
        setOptions(res.data.data);
        setEditOption(null);
        setEditingParagraph("");
      }
    );
  }

  function handleSubmitBulk(values: { input: string }) {
    const paragraphs = values.input.split("\n");

    let template_id = options[0].template_id;
    let templateString = template;
    let method = "POST";
    let body = paragraphs.map((name) => {
      return { name, template_id: options[0].template_id, user_id: user.id };
    });
    let route = USER_ROUTE + `/${user.id}/paragraph/bulk?template=${templateString}`;

    // Change request format if user selected a template.
    if (selectedTemplate) {
      template_id = templates[selectedTemplate].id;
      templateString = templates[selectedTemplate].name;
      method = "PUT";
      body = createUpdateParagraphs(options, paragraphs, template_id, user.id);
      route = USER_ROUTE + `/${user.id}/paragraph?template=${templateString}`;
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
            {options.map((_, i) => (
              <React.Fragment key={i}>
                <TableRow
                  columns={columns}
                  index={i}
                  items={options}
                  onClickEdit={() => {
                    setEditModal(true);
                    setEditingParagraph(options[i].name);
                    setEditOption(options[i]);
                  }}
                  onClickDelete={() => handleDelete(options[i].id!)}
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
              editingItem={options.map((op) => op.name).join("\n")}
              isLoading={isLoading}
              handleSubmit={handleSubmitBulk}
            />
          </Box>
        )}
      </Box>
    </Layout>
  );
};
