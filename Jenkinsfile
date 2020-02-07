#!/usr/bin/env groovy

def goleveranse
def config = [
    artifactId               : 'fiona',
    groupId                  : 'no.skatteetaten.aurora',
    deliveryBundleClassifier : "Doozerleveransepakke",
    scriptVersion            : 'v7',
    pipelineScript           : 'https://git.aurora.skead.no/scm/ao/aurora-pipeline-scripts.git',
    iqOrganizationName       : 'Team AOS',
    openShiftBaseImage       : 'turbo',
    openShiftBaseImageVersion: 'latest',
    goVersion                : 'Go 1.13',
    artifactPath             : 'bin/',
    sonarQube                : false,
    versionStrategy          : [
        [branch: 'master', versionHint: '1']
    ]
]

fileLoader.withGit(config.pipelineScript, config.scriptVersion) {
  goleveranse = fileLoader.load('templates/goleveranse')
}
