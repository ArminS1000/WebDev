const express = require("express");
const fs = require("fs");

const cors = require("cors");
const crypto = require("crypto");

var app = express();
app.use(cors());
app.use(express.json());

let fileText;
fs.readFile(
    "./text.txt",
    "utf8",
    (err, data) => (fileText = data.split("\n"))
);


app.post("/node/sha256", (req, res) => {
    const inputText = req.query.firstinput;

    if (typeof inputText == "undefined")
        return res.send("Error: empty inputs");

    if (inputText.length < 8)
        return res.send("Error: text has less than 8 charecters");

    const hash = crypto.createHash('sha256').update(inputText).toString().digest('hex');
    res.json({
        result: hash,
    });
});


app.get("/node/write", (req, res) => {
    let lineNumber  = req.query.input;

    if (typeof lineNumber == "undefined")
        return res.send("Error: empty inputs");

    fileText(async client => client.query('message').then(result => {
        if(result.length > 0)
            res.send(result);
        else 
            res.status(404).send("Error: text not found");
    }).catch(console.log));  
});

const port = 3000;
app.listen(port, () => {
    console.log(`listening on port ${port}...`);
});