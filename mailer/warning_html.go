package mailer

const warningHTML = `
<!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>

<head>
    <style>
        body {
            font-family: "Roboto", "Helvetica", "Arial", sans-serif;
        }

        table {
            font-family: arial, sans-serif;
            border-collapse: collapse;
            width: auto;
            border-collapse: collapse;
            border-spacing: 0;
            display: table;
            border: 1px solid #ccc;
        }

        tr {
            text-align: left;
            background-color: #fff;
            border-bottom: 1px solid #ddd;
        }

        tr:nth-child(event) {
            background-color: #fff;
        }

        tr:nth-child(odd) {
            background-color: #E7E9EB;
        }

        td {
            border: 1px solid #dddddd;
            padding: 0 15px 0 8px;
        }

        th {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
    </style>
</head>

<body>
    <p><b>Disk Warnings:</b></p>
    <table>
        <tr>
            <th>Mount on</th>
            <th>Use%</th>
        </tr>
        {{ range $i, $warning := .Warnings}}
        <tr>
            <td>{{ $warning.Device }}</td>
            <td>{{ $warning.Percent }}</td>
        </tr>
        {{ end }}
    </table>
    <p></p>
    <p>
        Kind Regards,
        <br />
        Disk Usage Warner
    </p>

</body>

</html>
`
