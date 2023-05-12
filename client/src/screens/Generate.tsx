import { Box, Button, FormLabel, useToast } from "@chakra-ui/react";
import React, { useCallback, useContext, useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import ReactSelect from "react-select";
import EditModal from "../components/EditModal";
import SentenceBox from "../components/SentenceBox";
import { USER_ROUTE } from "../constants";
import { UserContext } from "../context/UserContext";
import useFetch from "../hooks/useFetch";
import useLoginRequired from "../hooks/useLoginRequired";
import Layout from "../layout/Layout";
import {
  Content,
  Dictionary,
  DictionaryResponse,
  FinalizedContent,
  Sentence,
  Synonym,
  Template,
  Word,
} from "../types/general";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import { createUpdateSynonyms } from "../utils/createUpdateSynonyms";
import { extractTags } from "../utils/extractTags";
import { getRandomInt } from "../utils/getRandomInt";
import { transformContent } from "../utils/transformContent";
import { transformDictionary } from "../utils/transformDictionary";

const Generate: React.FC = () => {
  const { makeRequest, isLoading } = useFetch();
  const { user } = useContext(UserContext);
  const toast = useToast();
  const [selectedTemplate, setSelectedTemplate] = useState("");
  const [selectedWord, setSelectedWord] = useState<number>();
  const [templates, setTemplates] = useState<Template[]>([]);
  const [content, setContent] = useState<FinalizedContent[]>([]);
  const [dictionary, setDictionary] = useState<Dictionary>({});
  const [generatedContent, setGeneratedContent] = useState<Content[]>([]);
  const [editModal, setEditModal] = useState(false);
  const [editingSentences, setEditingSentences] = useState<Sentence[]>([]);
  const [editingSentencesParagraph, setEditingSentencesParagraph] =
    useState("");
  const [editSynonyms, setEditSynonyms] = useState<{
    word: string;
    synonyms: string[];
  }>({ word: "", synonyms: [""] });
  const [synonymModal, setSynonymModal] = useState(false);
  const [editingWord, setEditingWord] = useState<Word | null>();
  const [words, setWords] = useState<Word[]>([]);
  const [existingSynonyms, setExistingSynonyms] = useState<Synonym[]>([]);
  const [seeTagged, setSeeTagged] = useState(false);
  useLoginRequired();

  useEffect(() => {
    if (selectedWord) {
      makeRequest(
        {
          url: USER_ROUTE + `/${user.id}/synonym?word=${selectedWord}`,
        },
        (res) => {
          setExistingSynonyms(res.data.data);
        }
      );
    }
  }, [makeRequest, selectedWord, user.id]);

  useEffect(() => {
    // If no template has been selected, fetch templates.
    if (selectedTemplate.length === 0) {
      makeRequest(
        {
          url: USER_ROUTE + `/${user.id}/template`,
        },
        (res) => {
          setTemplates(res.data.data);
        }
      );
    } else {
      // Otherwise, pull the content associated with that template.
      makeRequest(
        {
          url: USER_ROUTE + `/${user.id}/content?template=${selectedTemplate}`,
        },
        (res) => {
          const initialContent: Content[] = res.data.data;
          setContent(transformContent(initialContent));
        }
      );

      // Pull dictionary
      makeRequest(
        {
          url: USER_ROUTE + `/${user.id}/dictionary`,
        },
        (res) => {
          const initialDictionary: DictionaryResponse[] = res.data.data;
          setDictionary(transformDictionary(initialDictionary));
        }
      );

      // On Change Template => set the content back to an empty array.
      setGeneratedContent([]);
    }
  }, [makeRequest, selectedTemplate, user.id]);

  // This function will use the Set Sentences to send a "bulk request" to the server.
  // The server will take care of the rest.
  const handleSubmit = useCallback(
    (values: { input: string }) => {
      const sentences = values.input.split("\n");
      const template_id = editingSentences[0].template_id;
      const paragraph_id = editingSentences[0].paragraph_id;
      const body = sentences.map((sentence) => {
        return {
          paragraph_id: paragraph_id,
          template_id: template_id,
          sentence,
        };
      });
      makeRequest(
        {
          url: USER_ROUTE + `/${user.id}/sentence/bulk?paragraph=${editingSentencesParagraph}`,
          method: "POST",
          data: body,
        },
        (res) => {
          setEditingSentences([]);
          setEditingSentencesParagraph("");
          setEditModal(false);
        }
      );
    },
    [makeRequest, editingSentences, editingSentencesParagraph, user.id]
  );

  const handleSynonyms = useCallback(
    (values: { input: string }) => {
      const synonyms = values.input.split("\n");

      let word_id = editingWord?.id;
      let wordString = editSynonyms.word;
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
        body = createUpdateSynonyms(existingSynonyms, synonyms, word_id, words[selectedWord]);
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
          setEditingSentences([]);
          setEditingSentencesParagraph("");
          setSynonymModal(false);
        }
      );
    },
    [
      editSynonyms.word,
      makeRequest,
      editingWord?.id,
      existingSynonyms,
      selectedWord,
      toast,
      words,
      user.id
    ]
  );

  const editSentence = (content: Content) => {
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/sentence?paragraph=${content.paragraph}`
      },
      (res) => {
        setEditingSentences(res.data.data);
      }
    );
    setEditingSentencesParagraph(content.paragraph);
    setEditModal(true);
  };

  const editWord = (word: string) => {
    const item = extractTags(word);
    setEditingWord(() => {
      return words.filter((word) => word.name === item.word)[0];
    });

    // Get the synonyms associated with the clicked word.
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/synonym?word=${item.word}`,
      },
      (res) => {
        setExistingSynonyms(res.data.data);
      }
    );

    // Get all of the user's words.
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/word`,
      },
      (res) => {
        setWords(res.data.data);
      }
    );

    const synonyms = dictionary[item.tag];

    if (synonyms) {
      setSynonymModal(true);
      setEditSynonyms({
        word: item.word,
        synonyms,
      });
    }
  };

  function selectContent(content: FinalizedContent[]): Content[] {
    let data: Content[] = [];
    for (let i = 0; i < content.length; i++) {
      const sentences =
        content[i].sentences[getRandomInt(content[i].sentences.length)];
      data.push({
        ...content[i],
        sentences,
      });
    }
    return data.sort((a, b) => a.order - b.order);
  }

  // Select Component for Words
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

  // Navigate to paragraph detail while on sentences modal.
  const navigateToParagraph = useCallback(() => {
    return (
      <Button variant={"outline"} colorScheme={"blue"}>
        <Link target="_blank" to={"/paragraph/" + editingSentencesParagraph}>
          Edit Sentence
        </Link>
      </Button>
    );
  }, [editingSentencesParagraph]);

  // Here's the "selecting template" logic
  const renderTopLevelOptions = useMemo(() => {
    return (
      <Box
        sx={{
          display: "flex",
          flexDirection: "row",
          my: 5,
          gap: 5,
          width: 500,
          justifyContent: "space-around",
          alignItems: "center",
        }}
      >
        <Button
          variant={"outline"}
          colorScheme={"blue"}
          onClick={() => setGeneratedContent(selectContent(content))}
        >
          Generate
        </Button>
        <Button
          variant={"outline"}
          colorScheme={"green"}
          onClick={() => setSeeTagged((prev) => !prev)}
        >
          See Tags
        </Button>
        <Box sx={{ width: 250 }}>
          <ReactSelect
            name={"select template"}
            placeholder={"select template"}
            value={{
              value: selectedTemplate,
              label: capitalizeFirstLetter(selectedTemplate),
            }}
            onChange={(e) => {
              setSelectedTemplate(String(e?.value));
            }}
            aria-label={"select template"}
            options={templates.map((op) => {
              return {
                value: op.name,
                label: capitalizeFirstLetter(op.name),
              };
            })}
          />
        </Box>
      </Box>
    );
  }, [content, selectedTemplate, templates]);

  // Modal for editing senteces
  const renderSentencesModal = useMemo(() => {
    return (
      <EditModal
        selectComponent={navigateToParagraph()}
        editModal={editModal}
        setEditModal={setEditModal}
        handleSubmit={handleSubmit}
        editingItem={editingSentences.map((s) => s.sentence).join("\n")}
        isLoading={isLoading}
      />
    );
  }, [
    editModal,
    handleSubmit,
    isLoading,
    editingSentences,
    navigateToParagraph,
  ]);

  // Modal for editing synonyms
  const renderSynonymsModal = useMemo(() => {
    return (
      <EditModal
        selectComponent={SelectChangeWord()}
        editModal={synonymModal}
        setEditModal={setSynonymModal}
        handleSubmit={handleSynonyms}
        editingItem={editSynonyms.synonyms.join("\n")}
        isLoading={isLoading}
      />
    );
  }, [
    editSynonyms.synonyms,
    handleSynonyms,
    isLoading,
    synonymModal,
    SelectChangeWord,
  ]);

  return (
    <Layout>
      {renderTopLevelOptions}
      <Box sx={{ height: "80vh", width: "60%" }}>
        <SentenceBox
          seeTagged={seeTagged}
          dict={dictionary}
          generatedContent={generatedContent}
          onClickWord={(word) => {
            if (!word.includes("@")) {
              editWord(word);
            }
          }}
          onClickSentence={(content) => {
            editSentence(content);
          }}
        />
      </Box>
      {editModal && renderSentencesModal}
      {synonymModal && renderSynonymsModal}
    </Layout>
  );
};

export default Generate;
