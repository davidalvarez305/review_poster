import { useMemo, useState } from "react";
import { Dictionary } from "../types/general";
import { spinContent } from "../utils/spinContent";

type WordBoxProps = React.DetailedHTMLProps<
  React.HTMLAttributes<HTMLSpanElement>,
  HTMLSpanElement
> & {
  word: string;
  dictionary: Dictionary;
  seeTagged: boolean;
};

export const WordBox: React.FC<WordBoxProps> = ({
  word,
  dictionary,
  seeTagged,
  ...props
}) => {
  const [isHovering, setIsHovering] = useState(false);

  // This prevents the spin content function from being called from changes in the parent 'Sentence' component.
  // For instance, if the Modals are called or if the sentence is hovered.
  
  const memoizedWord = useMemo(
    () => spinContent(word, dictionary),
    [word, dictionary]
  );

  let hoverStyles = {
    cursor: "pointer",
    color: "green",
    fontSize: 16,
    fontWeight: "bold",
  };

  let unmatchedStyles = { ...hoverStyles, color: "red" };

  return (
    <span
      {...props}
      onMouseEnter={() => {
        setIsHovering(true);
      }}
      onMouseLeave={() => {
        setIsHovering(false);
      }}
      style={
        isHovering
          ? hoverStyles
          : word.includes("(#")
          ? undefined
          : unmatchedStyles
      }
    >
      {seeTagged ? word : memoizedWord}
    </span>
  );
};
