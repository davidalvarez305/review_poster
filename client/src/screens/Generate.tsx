import { Box, Button } from "@chakra-ui/react";
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
import {
  Dictionary,
  Paragraph,
  Sentence,
  SpunContent,
  Word,
} from "../types/general";
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
  const { templates, getUserTemplates } = useTemplatesController();
  const { words, getUserWords } = useWordsController();
  const {
    joinedSentences,
    updateUserParagraphSentencesByTemplate,
    GetUserJoinedSentencesByParagraph,
  } = useSentencesController();
  const { updateUserSynonymsByWord, getUserSynonymsByWord } =
    useSynonymsController();
  useLoginRequired();

  useEffect(() => {
    // If no template has been selected, fetch templates.
    if (selectedTemplate.length === 0) {
      getUserTemplates();
    } else {
      GetUserJoinedSentencesByParagraph(selectedTemplate);

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
    getUserTemplates,
    GetUserJoinedSentencesByParagraph,
  ]);

  const handleSubmit = useCallback(
    (values: { input: string }, paragraph: Paragraph) => {
      /* const template_id = editingSentences[0].paragraph!.template_id;
      const paragraph_id = editingSentences[0].paragraph_id; */
      updateUserParagraphSentencesByTemplate(
        { ...values },
        paragraph.name,
        selectedTemplate
      );

      setEditingSentences([]);
      setEditingSentencesParagraph("");
      setEditModal(false);
    },
    [updateUserParagraphSentencesByTemplate, selectedTemplate]
  );

  const handleSynonyms = useCallback(
    (values: { input: string }) => {
      if (editingWord) {
        updateUserSynonymsByWord({ ...values }, editingWord.name);
      }

      setEditingSentences([]);
      setEditingSentencesParagraph("");
      setSynonymModal(false);
    },
    [updateUserSynonymsByWord, editingWord]
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

  function selectContent(paragraphs: Paragraph[]): SpunContent[] {
    let data: SpunContent[] = [];

    for (const paragraph of paragraphs) {
      if (!paragraph.sentences || paragraph.sentences.length === 0) continue;

      const sentence = paragraph.sentences[getRandomInt(paragraph.sentences.length)];

      data.push({
        paragraph: paragraph.name,
        sentence: sentence.sentence,
        order: 0,
      });
    }
    return data.sort((a, b) => a.order - b.order);
  }

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
          onClick={() => setGeneratedContent(selectContent(joinedSentences))}
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
  }, [joinedSentences, selectedTemplate, templates]);

  // Modal for editing senteces
  const renderSentencesModal = useMemo(() => {
    let paragraph: Paragraph | null = null;
    if (editingSentences.length > 0)
      paragraph = editingSentences.filter((s) => s.sentence)[0].paragraph;
    return (
      <EditModal
        selectComponent={navigateToParagraph()}
        editModal={editModal}
        setEditModal={setEditModal}
        handleSubmit={(values) => {
          if (paragraph) handleSubmit(values, paragraph);
        }}
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
        editModal={synonymModal}
        setEditModal={setSynonymModal}
        handleSubmit={handleSynonyms}
        editingItem={editSynonyms.synonyms.join("\n")}
        isLoading={isLoading}
      />
    );
  }, [editSynonyms.synonyms, handleSynonyms, isLoading, synonymModal]);

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
            let sentence: Sentence | null = null;

            for (let i = 0; i < joinedSentences.length; i++) {
              if (!joinedSentences[i].sentences) continue;

              for (let n = 0; n < joinedSentences[i].sentences!.length; n++) {
                if (joinedSentences[i].sentences![n].sentence === content.sentence) {
                  sentence = joinedSentences[i].sentences![n];
                  break;
                }
              }
            }

            if (sentence) editSentence(sentence);
          }}
        />
      </Box>
      {editModal && renderSentencesModal}
      {synonymModal && renderSynonymsModal}
    </Layout>
  );
};

export default Generate;
