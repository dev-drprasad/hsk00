import React from "react";
import "./GameList.scss";

const getClassName = (g) => {
  let className = "";
  if (!g.hsk) className += "unsaved ";
  if (g.deleted) className += "mark-deleted";
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
              <span role="img" aria-label="delete icon">
                {g.deleted ? "✅" : "❌"}
              </span>
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default GameList;
