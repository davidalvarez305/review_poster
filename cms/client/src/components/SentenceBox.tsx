import React, { useEffect, useState } from "react";
import { Content, Dictionary } from "../types/general";
import { WordBox } from "./WordBox";

interface Props {
  dict: Dictionary;
  generatedContent: Content[];
  onClickSentence: (content: Content) => void;
  onClickWord: (word: string) => void;
  seeTagged: boolean;
}

const SentenceBox: React.FC<Props> = ({
  dict,
  generatedContent,
  onClickSentence,
  onClickWord,
  seeTagged,
}) => {
  const [content, setContent] = useState<Content[]>([]);
  const [hoveringWord, setHoveringWord] = useState(false);

  useEffect(() => {
    const vals = generatedContent.map((sentence) => {
      let values: string[] = [];
      let words = sentence.sentences.split(/ /);

      if (words) {
        for (let n = 0; n < words.length; n++) {
          if (!Array.isArray(words[n])) {
            values.push(words[n]);
          } else {
            values = [...values, ...words[n]];
          }
        }
      }
      return { ...sentence, sentences: values.join(" ") };
    });
    setContent(vals);
  }, [generatedContent, hoveringWord]);

  return (
    <>
      {content.map((p) => {
        const sentences = p.sentences.split(/[ .?!]/);
        return (
          <div
            onClick={() => {
              // Only trigger sentence modal when words are not being clicked/hovered.
              if (!hoveringWord) {
                onClickSentence(p);
              }
            }}
            style={{
              fontFamily: "Georgia",
              fontSize: 18,
              marginTop: 20,
            }}
            key={p.sentences}
          >
            {sentences.map((str, index) => (
              <span
                key={index}
                onMouseEnter={() => {
                  setHoveringWord(true);
                }}
                onMouseLeave={() => {
                  setHoveringWord(false);
                }}
              >
                <WordBox
                  word={str}
                  dictionary={dict}
                  onClick={() => {
                    onClickWord(str);
                  }}
                  id={str}
                  seeTagged={seeTagged}
                />
                {" "}
              </span>
            ))}
          </div>
        );
      })}
    </>
  );
};

export default SentenceBox;
