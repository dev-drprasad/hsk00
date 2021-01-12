import React, { useState } from "react";
import "./App.scss";
import Alert from "./Alert";
import GitHubIcon from "./GitHubIcon";

const initialState = {
  rootDir: "",
  newGames: [],
  categoryID: -1,
  errors: { rootDir: "" },
};

function App() {
  const [state, setState] = useState(initialState);

  const handleSelectRootClick = () => {
    window.backend.Runtime.SelectRootDir().then((selectedDir) => {
      setState((s) => ({ ...s, rootDir: selectedDir }));
    });
  };

  const handleSelectGamesClick = (e) => {
    window.backend.Runtime.SelectGames().then((selectedGames) => {
      console.log("selectedGames :>> ", selectedGames);
      setState((s) => ({ ...s, newGames: selectedGames }));
    });
  };

  const handleRootDirChange = (e) => {
    setState((s) => ({ ...s, rootDir: e.target.value }));
  };

  const handleCategoryChange = (e) => {
    setState((s) => ({ ...s, categoryID: e.target.value }));
  };

  const handleSubmit = () => {
    console.log(
      "rootDir, categoryID, newGames :>> ",
      rootDir,
      categoryID,
      newGames
    );
  };

  const { rootDir, newGames, categoryID, errors } = state;
  return (
    <React.Fragment>
      <div className="App">
        <div className="FormItem">
          <div className="Label" htmlFor="rootDir">
            Choose root path :
          </div>
          <div className="group RootDirGroup">
            <input
              className="FormControl Input"
              name="rootDir"
              placeholder="Choose root path (SD Card)"
              onChange={handleRootDirChange}
              value={rootDir}
            />
            <button className="FormControl btn" onClick={handleSelectRootClick}>
              Choose
            </button>
          </div>
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <div className="FormItem">
          <div className="Label" htmlFor="rootDir">
            Select game category :
          </div>
          <div>
            <select
              className="FormControl Select CategorySelect"
              name="categoryID"
              onChange={handleCategoryChange}
            >
              <option value={-1}>----------</option>
              <option value={0}>0. Action Games</option>
              <option value={1}>1. Shoot Games</option>
              <option value={2}>2. Sport Games</option>
              <option value={3}>3. Fight Games</option>
              <option value={4}>4. Racing Games</option>
              <option value={5}>5. Puzzle Games</option>
            </select>
          </div>
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <div className="FormItem">
          <div className="Label" htmlFor="rootDir">
            Games to add :
          </div>
          <div className="ListBox">
            <ul role="listbox">
              {newGames.map((g) => (
                <li key={g}>{g}</li>
              ))}
            </ul>
            <button
              className="FormControl btn"
              onClick={handleSelectGamesClick}
            >
              Select games
            </button>
          </div>
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <Alert type="danger" message="Hooray! 14 games added!" />
        <div className="FormItem SubmitButtonWrapper">
          <button
            className="FormControl SubmitButton btn btn-lg btn-primary"
            disabled={newGames.length === 0}
            onClick={handleSubmit}
          >
            Add {newGames.length} games
          </button>
        </div>
      </div>
      <a className="github-link" href="https://github.com/dev-drprasad/hsk00">
        <GitHubIcon />
      </a>
    </React.Fragment>
  );
}

export default App;
