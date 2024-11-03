import React from 'react';
import { AuthProvider } from './context/AuthContext';
import ItemList from './ItemList';

function App() {
    return (
        <AuthProvider>
            <div className="App">
                <ItemList />
            </div>
        </AuthProvider>
    );
}

export default App;
