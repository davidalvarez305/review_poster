import React, { useEffect, useState } from "react";
import { Dictionary, SpunContent } from "../types/general";
import { WordBox } from "./WordBox";

interface Props {
  dict: Dictionary;
  generatedContent: SpunContent[];
  onClickSentence: (content: SpunContent) => void;
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
  const [content, setContent] = useState<SpunContent[]>([]);
  const [hoveringWord, setHoveringWord] = useState(false);

  useEffect(() => {
    const vals = generatedContent.map((content) => {
      let values: string[] = [];
      let words = content.sentence.split(/ /);

      if (words) {
        for (let n = 0; n < words.length; n++) {
          if (!Array.isArray(words[n])) {
            values.push(words[n]);
          } else {
            values = [...values, ...words[n]];
          }
        }
      }
      return { ...content, sentence: values.join(" ") };
    });
    setContent(vals);
  }, [generatedContent, hoveringWord]);

  return (
    <>
      {content.map((spunContent) => {
        const sentences = spunContent.sentence.split(/[ .?!]/);
        return (
          <div
            onClick={() => {
              // Only trigger sentence modal when words are not being clicked/hovered.
              if (!hoveringWord) {
                onClickSentence(spunContent);
              }
            }}
            style={{
              fontFamily: "Georgia",
              fontSize: 18,
              marginTop: 20,
            }}
            key={spunContent.sentence}
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
