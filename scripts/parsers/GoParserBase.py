# encoding: utf-8

import sys
if sys.version_info[1] > 5:
    from typing import TextIO
else:
    from typing.io import TextIO

import antlr4
from antlr4 import Lexer
from antlr4 import Parser
from antlr4 import TokenStream

from .GoLexer import GoLexer


# GoParserBase implementation
class GoParserBase(Parser):

    def __init__(self, input: TokenStream, output: TextIO = sys.stdout):
        super().__init__(input, output)

    # Returns `True` if on the current index of the parser's
    # token stream a token exists on the `HIDDEN` channel which
    # either is a line terminator, or is a multi line comment that
    # contains a line terminator.
    def lineTerminatorAhead(self) -> bool:
        offset = 1
        possibleIndexEosToken = self.getCurrentToken().tokenIndex - offset
        if possibleIndexEosToken == -1:
            return True

        ahead = self.getTokenStream().get(possibleIndexEosToken)
        while ahead.channel == Lexer.HIDDEN:
            if ahead.type == GoLexer.TERMINATOR:
                return True
            if ahead.type == GoLexer.WS:
                offset += 1
                possibleIndexEosToken = self.getCurrentToken().tokenIndex - offset
                ahead = self.getTokenStream().get(possibleIndexEosToken)
            if ahead.type == GoLexer.COMMENT or ahead.type == GoLexer.LINE_COMMENT:
                if '\r' in ahead.text or '\n' in ahead.text:
                    return True
                else:
                    offset += 1
                    possibleIndexEosToken = self.getCurrentToken().tokenIndex - offset
                    ahead = self.getCurrentToken().get(possibleIndexEosToken)
        return False

    # Returns `True` if no line terminator exists between the specified
    # token offset and the prior one on the `HIDDEN` channel.
    def noTerminatorBetween(self, tokenOffset: int) -> bool:
        stream = self.getTokenStream()
        tokens = stream.getHiddenTokensToLeft(stream.LT(tokenOffset).tokenIndex, -1)
        if tokens is None:
            return True
        for token in tokens:
            if '\n' in token.text:
                return False
        return True

    # Returns `True` if no line terminator exists after any encountered
    # parameters beyond the specified token offset and the next on the
    # `HIDDEN` channel.
    def noTerminatorAfterParams(self, tokenOffset: int) -> bool:
        leftParams = 1
        rightParams = 0
        stream = self.getTokenStream()
        if stream.LT(tokenOffset).type == GoLexer.L_PAREN:
            while leftParams != rightParams:
                tokenOffset += 1
                tokenType = stream.LT(tokenOffset).type
                if tokenType == GoLexer.L_PAREN:
                    leftParams += 1
                elif tokenType == GoLexer.R_PAREN:
                    rightParams += 1
            tokenOffset += 1
            return self.noTerminatorBetween(tokenOffset)
        return True

    def checkPreviousTokenText(self, text: str) -> bool:
        stream = self.getTokenStream()
        return stream.LT(1).text == text
