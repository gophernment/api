FROM scratch

ENV PORT 80
EXPOSE $PORT

COPY api /
CMD ["/api"]