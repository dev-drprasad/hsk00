import React from "react";
import Emoji from "./Emoji";
import "./GameList.scss";

const getClassName = (g) => {
  let className = "";
  if (!g.hsk) className += " unsaved ";
  if (g.deleted) className += " mark-deleted ";
  return className;
};

function GameList({ games, onToggleDelete }) {
  return (
    <div className="ListBox">
      <ul role="listbox">
        {games.map((g, i) => (
          <li key={`${g.filename}${g.srcPath}${g.id}`} className={getClassName(g)}>
            {g.name}
            <button className="delete-btn btn btn-sm" onClick={onToggleDelete(i)}>
              <span className="emoji" role="img" aria-label="delete icon"></span>
              <Emoji ariaLabel="delete icon">{g.deleted ? "✅" : "❌"}</Emoji>
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default GameList;
