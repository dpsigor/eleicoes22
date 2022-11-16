import re
import json
import argparse
import sys
from ast import literal_eval

INDENT = 4
sys.stdout.reconfigure(encoding="utf-8")
sys.stdin.reconfigure(encoding="utf-8")

def toObject(lines):
    def __toObject(i, level):
        def getlevel(line):
            return int((len(line) - len(lstrip(line)))/4)
        def lstrip(line):
            return re.sub(r"^[\s|\.]+", "" , line)

        obj = {}
        repeat = False
        while (i < len(lines) and getlevel(lines[i]) >= level):
            line_adjusted = lstrip(lines[i])
            level = getlevel(lines[i])

            if '] <==' in line_adjusted: # fim de lista
                i = i + 1
            elif line_adjusted.endswith(':'): #objetos
                (obj[line_adjusted.replace(':','')], i, _) = __toObject(i+1, level+1)
            elif "=" in line_adjusted: # literais
                [k, v] = [l.strip() for l in line_adjusted.split("=")]
                if(k in obj):
                    repeat = True
                    break
                else:
                    i = i + 1
                    try:
                        obj[k] = literal_eval(v)
                    except:
                        obj[k] = v
            elif line_adjusted.endswith(': ['): # listas
                k = line_adjusted.replace(': [','')
                obj[k] = []
                __repeat = True
                i = i + 1
                while __repeat:
                    (__obj, i, __repeat) = __toObject(i, level+1)
                    obj[k].append(__obj)
        return (obj, i, repeat)
    return __toObject(0, 0)[0]

parser = argparse.ArgumentParser()
parser.add_argument('-i', '--input', help='Arquivo de entrada', required=False)
parser.add_argument('-o', '--output', help='Arquivo de saída', required=False)
args = parser.parse_args()

if args.input:
    with open(args.input, encoding='utf-8') as file:
        lines = file.readlines()
else:
    lines = sys.stdin.readlines()
dict = toObject([line.rstrip() for line in lines])

if args.output:
    with open(args.output, 'w', encoding='utf-8') as file:
        json.dump(dict, file, ensure_ascii=False, indent=INDENT)
else:
    json.dump(dict, sys.stdout, ensure_ascii=False, indent=INDENT)