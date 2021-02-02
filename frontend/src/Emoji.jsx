import React from "react";

function Emoji({ children, ariaLabel }) {
  return (
    <span className="emoji" role="img" aria-label={ariaLabel}>
      {children}
    </span>
  );
}

export default Emoji;
