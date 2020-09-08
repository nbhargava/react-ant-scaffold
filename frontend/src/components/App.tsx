import * as React from 'react';
import { Button, message } from 'antd';

import '../style/App.css';

class App extends React.Component {
  public render() {
    const m = () => message.success('javascript works');

    return (
      <div className="main-content">
        <div className="center">
          <Button onClick={m}>Sup</Button>
        </div>
      </div>
    );
  }
}

export default App;
