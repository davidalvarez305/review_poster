import { Box, Button, FormLabel } from "@chakra-ui/react";
import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { Link } from "react-router-dom";
import ReactSelect from "react-select";
import EditModal from "../components/EditModal";
import SentenceBox from "../components/SentenceBox";
import { USER_ROUTE } from "../constants";
import { UserContext } from "../context/UserContext";
import useFetch from "../hooks/useFetch";
import useLoginRequired from "../hooks/useLoginRequired";
import Layout from "../layout/Layout";
import { Dictionary, Sentence, SpunContent, Word } from "../types/general";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import { extractTags } from "../utils/extractTags";
import { getRandomInt } from "../utils/getRandomInt";
import { createDictionaryFactory } from "../utils/createDictionaryFactory";
import useTemplatesController from "../hooks/useTemplatesController";
import useWordsController from "../hooks/useWordsController";
import useSentencesController from "../hooks/useSentencesController";
import useSynonymsController from "../hooks/useSynonymsController";

const Generate: React.FC = () => {
  const { makeRequest, isLoading } = useFetch();
  const { user } = useContext(UserContext);
  const [selectedTemplate, setSelectedTemplate] = useState("");
  const [selectedWord, setSelectedWord] = useState<number>();
  const [dictionary, setDictionary] = useState<Dictionary>({});
  const [generatedContent, setGeneratedContent] = useState<SpunContent[]>([]);
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
  const [seeTagged, setSeeTagged] = useState(false);
  const { templates, getTemplates } = useTemplatesController();
  const { words, getUserWords } = useWordsController();
  const { sentences, getSentencesByTemplate, bulkUpdateSentences } =
    useSentencesController();
  const { updateUserSynonymsByWord, updateSynonyms, getUserSynonymsByWord } =
    useSynonymsController();
  useLoginRequired();

  useEffect(() => {
    if (selectedWord) {
      getUserSynonymsByWord(words[selectedWord].name);
    }
  }, [getUserSynonymsByWord, selectedWord, user.id, words]);

  useEffect(() => {
    // If no template has been selected, fetch templates.
    if (selectedTemplate.length === 0) {
      getTemplates();
    } else {
      getSentencesByTemplate(selectedTemplate);

      // Pull dictionary
      makeRequest(
        {
          url: USER_ROUTE + `/${user.id}/word`,
        },
        (res) => {
          const initialDictionary: Word[] = res.data.data;
          setDictionary(createDictionaryFactory(initialDictionary));
        }
      );

      // On Change Template => set the content back to an empty array.
      setGeneratedContent([]);
    }
  }, [
    makeRequest,
    selectedTemplate,
    user.id,
    getSentencesByTemplate,
    getTemplates,
  ]);

  const handleSubmit = useCallback(
    (values: { input: string }) => {
      /* const template_id = editingSentences[0].paragraph!.template_id;
      const paragraph_id = editingSentences[0].paragraph_id; */
      bulkUpdateSentences({ ...values });

      setEditingSentences([]);
      setEditingSentencesParagraph("");
      setEditModal(false);
    },
    [bulkUpdateSentences]
  );

  const handleSynonyms = useCallback(
    (values: { input: string }) => {
      if (editingWord) {
        updateUserSynonymsByWord({ ...values }, editingWord.name);
      }
      // Change request format if user selected a word.
      if (selectedWord) {
        updateSynonyms(
          { ...values },
          words[selectedWord].id!,
          words[selectedWord].name
        );
      }

      setEditingSentences([]);
      setEditingSentencesParagraph("");
      setSynonymModal(false);
    },
    [updateUserSynonymsByWord, editingWord, selectedWord, updateSynonyms, words]
  );

  const editSentence = (sentence: Sentence) => {
    if (sentence.paragraph) {
      makeRequest(
        {
          url:
            USER_ROUTE +
            `/${user.id}/sentence?paragraph=${sentence.paragraph.name}`,
        },
        (res) => {
          setEditingSentences(res.data.data);
        }
      );
      setEditingSentencesParagraph(sentence.paragraph.name);
      setEditModal(true);
    }
  };

  const editWord = (word: string) => {
    const item = extractTags(word);
    setEditingWord(() => {
      return words.filter((word) => word.name === item.word)[0];
    });

    getUserSynonymsByWord(word);

    getUserWords();

    const synonyms = dictionary[item.tag];

    if (synonyms) {
      setSynonymModal(true);
      setEditSynonyms({
        word: item.word,
        synonyms,
      });
    }
  };

  function selectContent(sentences: Sentence[]): SpunContent[] {
    let data: SpunContent[] = [];

    let adjustedSentences: { [key: string]: string[] } = {};

    for (let i = 0; i < sentences.length; i++) {
      if (!sentences[i].paragraph) continue;

      const key = sentences[i].paragraph!.name;
      if (adjustedSentences[key]) {
        adjustedSentences[key] = [
          ...adjustedSentences[key],
          sentences[i].sentence,
        ];
      } else {
        adjustedSentences[key] = [sentences[i].sentence];
      }
    }

    for (const [key, value] of Object.entries(adjustedSentences)) {
      const sentence = value[getRandomInt(value.length)];

      data.push({
        paragraph: key,
        sentence: sentence,
        order: 0,
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
          onClick={() => setGeneratedContent(selectContent(sentences))}
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
  }, [sentences, selectedTemplate, templates]);

  // Modal for editing senteces
  const renderSentencesModal = useMemo(() => {
    console.log("editingSentences: ", editingSentences);
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
    const word = selectedWord ? words[selectedWord].name : "";
    return (
      <EditModal
        selectComponent={SelectChangeWord()}
        editModal={synonymModal}
        setEditModal={setSynonymModal}
        handleSubmit={handleSynonyms}
        editingItem={editSynonyms.synonyms.join("\n")}
        isLoading={isLoading}
        selectedWord={word}
      />
    );
  }, [
    editSynonyms.synonyms,
    handleSynonyms,
    isLoading,
    synonymModal,
    SelectChangeWord,
    selectedWord,
    words,
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
            editSentence(
              sentences.filter(
                (sentence) => sentence.sentence === content.sentence
              )[0]
            );
          }}
        />
      </Box>
      {editModal && renderSentencesModal}
      {synonymModal && renderSynonymsModal}
    </Layout>
  );
};

export default Generate;
